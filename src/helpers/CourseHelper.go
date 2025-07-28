package helpers

import (
	"errors"
	"fmt"
	"goscraper/src/types"
	"goscraper/src/utils"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
)

type CoursePage struct {
	cookie string
}

func NewCoursePage(cookie string) *CoursePage {
	return &CoursePage{
		cookie: cookie,
	}
}

func (c *CoursePage) getUrl(currentDate time.Time) string {
	currentMonth := currentDate.Month()
	currentYear := currentDate.Year()

	var academicYearStart, academicYearEnd int

	if currentMonth >= 8 && currentMonth <= 12 {
		academicYearStart = currentYear - 1
		academicYearEnd = currentYear
	} else {
		academicYearStart = currentYear - 2
		academicYearEnd = currentYear - 1
	}

	url := fmt.Sprintf("https://academia.srmist.edu.in/srm_university/academia-academic-services/page/My_Time_Table_%d_%d", academicYearStart, academicYearEnd%100)
	return url
}

func (c *CoursePage) GetPage() (string, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(c.getUrl(time.Now()))
	req.Header.SetMethod("GET")
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("cookie", utils.ExtractCookies(c.cookie))
	req.Header.Set("Referer", "https://academia.srmist.edu.in/")
	req.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	req.Header.Set("Cache-Control", "private, max-age=120, must-revalidate")

	if err := fasthttp.Do(req, resp); err != nil {
		return "", fmt.Errorf("failed to fetch page: %v", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return "", fmt.Errorf("server returned status %d", resp.StatusCode())
	}

	data := string(resp.Body())
	parts := strings.Split(data, ".sanitize('")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid response format")
	}

	htmlHex := strings.Split(parts[1], "')")[0]
	return utils.ConvertHexToHTML(htmlHex), nil
}

func (c *CoursePage) GetCourses() (*types.CourseResponse, error) {
	page, err := c.GetPage()

	if err != nil {
		return &types.CourseResponse{
			Status: 500,
			Error:  err.Error(),
		}, err
	}

	re := regexp.MustCompile(`RA2\d{12}`)
	regNumber := re.FindString(page)

	htmlParts := strings.Split(page, `<table cellspacing="1" cellpadding="1" border="1" align="center" style="width:900px!important;" class="course_tbl">`)
	if len(htmlParts) < 2 {
		return &types.CourseResponse{
			Status: 500,
			Error:  "failed to find course table in the page",
		}, errors.New("failed to find course table in the page")
	}
	html := htmlParts[1]
	html = strings.Split(html, "</table>")[0]

	html = `<table style="font-size :16px;" border="1" align="center" cellpadding="1" cellspacing="1" bgcolor="#FAFAD2"><tbody>` + html + "</tbody></table>"

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return &types.CourseResponse{
			Status: 500,
			Error:  fmt.Sprintf("failed to parse HTML: %v", err),
		}, nil
	}

	var courses []types.Course
	rows := doc.Find("tr")

	rows.Each(func(i int, row *goquery.Selection) {
		if i == 0 {
			return
		}

		cells := row.Find("td")
		if cells.Length() > 0 {
			course := c.parseCourseRow(cells)
			if course != nil {
				courses = append(courses, *course)
			}
		}
	})

	return &types.CourseResponse{
		RegNumber: regNumber,
		Courses:   courses,
	}, nil
}

func (c *CoursePage) parseCourseRow(cells *goquery.Selection) *types.Course {
	if cells.Length() < 11 {
		return nil
	}

	getText := func(index int) string {
		return strings.TrimSpace(cells.Eq(index).Text())
	}

	code := getText(1)
	title := getText(2)
	credit := getText(3)
	category := getText(4)
	courseCategory := getText(5)
	courseType := getText(6)
	faculty := getText(7)
	slot := getText(8)
	room := getText(9)
	academicYear := getText(10)

	if credit == "" {
		credit = "N/A"
	}
	if courseType == "" {
		courseType = "N/A"
	}
	if faculty == "" {
		faculty = "N/A"
	}
	if room == "" {
		room = "N/A"
	} else {
		room = strings.ToUpper(room[:1]) + room[1:]
	}
	slot = strings.TrimSuffix(slot, "-")

	return &types.Course{
		Code:           code,
		Title:          strings.Split(title, " \\u2013")[0],
		Credit:         credit,
		Category:       category,
		CourseCategory: courseCategory,
		Type:           courseType,
		SlotType:       c.getSlotType(slot),
		Faculty:        faculty,
		Slot:           slot,
		Room:           room,
		AcademicYear:   academicYear,
	}
}

func (c *CoursePage) getSlotType(slot string) string {
	if strings.Contains(slot, "P") {
		return "Practical"
	}
	return "Theory"
}

func getYear(registrationNumber string) int {
	yearString := registrationNumber[2:4]
	currentYear := time.Now().Year()
	currentMonth := time.Now().Month()
	currentYearLastTwoDigits := currentYear % 100

	academicYearLastTwoDigits := utils.ParseInt(yearString)

	academicYear := currentYearLastTwoDigits
	if currentMonth >= 7 {
		academicYear++
	}

	studentYear := academicYear - academicYearLastTwoDigits

	if academicYearLastTwoDigits > currentYearLastTwoDigits {
		studentYear--
	}

	return studentYear
}
