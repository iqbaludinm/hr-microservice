package helper

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// This function is used for rendering a pdf from a template.
// The template is located in 'templates' folder, that is embedded in the executable of this apps.
// The flow of this function is:
//
// 1. Render the html from the template.
//
// 2. Create a buffer to add the html to the request form-data.
//
// 3. Create the http request to gotenberg service.
//
// 4. Send the request to gotenberg service.
//
// 5. Copy the rendered pdf from response Body to 'pdfResult' variable.
//
// 6. Return the 'pdfResult' variable.
func RenderPDF(templateFS embed.FS, templateName string, landscape bool, data interface{}) (bytes.Buffer, error) {
	// Define the functions
	funcMap := template.FuncMap{
		"add": func(x int) int { return x + 1 },
		"add_values": func(x int, yString string) int {
			y, _ := strconv.Atoi(yString)
			return x + y
		},
		"timestamp_to_date": func(x time.Time) string { return x.Format("2006-01-02") },
		"variance_background_color": func(variance int) template.HTMLAttr {
			var style string
			if variance < 0 {
				style = "style='background-color: #C6EFCD; color: #549E57;'"
			} else {
				style = "style='background-color: #FFC7CD; color: #CB5F65;'"
			}
			return template.HTMLAttr(style)
		},
		"forecasting_get_model_data": func(data string) string {
			return strings.Split(data, ",")[0]
		},
		"forecasting_get_component_description_data": func(data string) string {
			return strings.Split(data, ",")[1]
		},
	}
	// render html from template
	parsedTemplate, err := template.New(templateName).Funcs(funcMap).ParseFS(templateFS, fmt.Sprintf("templates/%s", templateName))
	if err != nil {
		return bytes.Buffer{}, err
	}

	var html bytes.Buffer
	err = parsedTemplate.Execute(&html, data)
	if err != nil {
		return bytes.Buffer{}, err
	}

	// Create a buffer to add the html to the request form-data
	var formData bytes.Buffer
	w := multipart.NewWriter(&formData)

	fw, err := w.CreateFormFile("index.html", "index.html")
	if err != nil {
		return bytes.Buffer{}, err
	}
	if _, err = io.Copy(fw, &html); err != nil {
		return bytes.Buffer{}, err
	}

	// set landscape or portrait
	if landscape {
		fieldWriter, err := w.CreateFormField("landscape")
		if err != nil {
			return bytes.Buffer{}, err
		}
		_, err = fieldWriter.Write([]byte("true"))
		if err != nil {
			return bytes.Buffer{}, err
		}
	}

	if err = w.Close(); err != nil {
		return bytes.Buffer{}, err
	}

	// Create the http request to gotenberg service
	req, err := http.NewRequest("POST", gotenbergEndpoint, &formData)
	if err != nil {
		return bytes.Buffer{}, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Send the request to gotenberg service
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return bytes.Buffer{}, err
	}
	if resp.StatusCode != 200 {
		return bytes.Buffer{}, fmt.Errorf("gotenberg service return status: %s \n %s", resp.Status, resp.Body)
	}

	// copy the rendered pdf from response Body to 'pdfResult' variable
	var pdfResult bytes.Buffer
	if _, err := io.Copy(&pdfResult, resp.Body); err != nil {
		return bytes.Buffer{}, err
	}
	defer resp.Body.Close()

	return pdfResult, nil
}
