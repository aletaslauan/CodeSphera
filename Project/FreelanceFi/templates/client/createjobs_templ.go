// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"freelancefi/db"
	"strings"
)

func CreateJobPage(username string, categories []db.Category) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><title>Create Job - FreelanceFi</title><link href=\"https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css\" rel=\"stylesheet\"><link href=\"https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.5/font/bootstrap-icons.css\" rel=\"stylesheet\"></head><body class=\"bg-light\"><!-- Top Navbar --><nav class=\"navbar navbar-expand-lg navbar-dark bg-dark shadow-sm border-bottom px-4 py-2 mb-0\"><div class=\"container-fluid\"><a class=\"navbar-brand fw-bold text-white\" href=\"/home\" style=\"margin-right: 60px; font-size: 1.75rem;\">FreelanceFi</a><div class=\"d-flex align-items-center gap-3 ms-auto\"><a href=\"/messages\" class=\"text-white\"><i class=\"bi bi-envelope fs-5\"></i></a> <a href=\"/saved\" class=\"text-white\"><i class=\"bi bi-bookmark fs-5\"></i></a> <a href=\"/notifications\" class=\"text-white\"><i class=\"bi bi-bell fs-5\"></i></a><div class=\"dropdown position-relative\" id=\"profileDropdownWrapper\"><button class=\"btn btn-outline-light rounded-circle text-uppercase fw-bold\" type=\"button\" id=\"profileDropdown\" style=\"width: 40px; height: 40px;\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if len(username) > 0 {
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(strings.ToUpper(string(username[0])))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/client/createjobs.templ`, Line: 30, Col: 66}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs("?")
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/client/createjobs.templ`, Line: 32, Col: 33}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</button><ul class=\"dropdown-menu position-absolute end-0 mt-2\" aria-labelledby=\"profileDropdown\" style=\"min-width: 180px;\"><li><a class=\"dropdown-item\" href=\"/profile\"><i class=\"bi bi-person me-2\"></i>Profile</a></li><li><a class=\"dropdown-item\" href=\"/settings\"><i class=\"bi bi-gear me-2\"></i>Settings</a></li><li><hr class=\"dropdown-divider\"></li><li><a class=\"dropdown-item text-danger\" href=\"/logout\"><i class=\"bi bi-box-arrow-right me-2\"></i>Log out</a></li></ul></div></div></div></nav><!-- Layout: Sidebar + Main Content --><div class=\"d-flex\" style=\"height: calc(100vh - 64px);\"><!-- Sidebar --><div class=\"bg-dark text-white d-flex flex-column p-3\" style=\"width: 250px; flex-shrink: 0;\"><ul class=\"nav flex-column mb-auto\"><li class=\"nav-item mb-2\"><a class=\"nav-link text-white\" href=\"/home\">Home</a></li><li class=\"nav-item mb-2\"><a class=\"nav-link text-white\" href=\"/jobspage\">Jobs</a></li><li class=\"nav-item mb-2\"><a class=\"nav-link text-white active\" href=\"/createjobs\">Post a Job</a></li><li class=\"nav-item mb-2\"><a class=\"nav-link text-white\" href=\"/mywork\">My Work</a></li><li class=\"nav-item mb-2\"><a class=\"nav-link text-white\" href=\"/finance\">Finance</a></li></ul><div class=\"mt-4\"><a class=\"btn btn-secondary w-100\" href=\"/logout\">Logout</a></div></div><!-- Main Content --><div class=\"flex-grow-1 overflow-auto p-4\"><h2 class=\"mb-4\">Post a New Job</h2><div class=\"card shadow-sm\"><div class=\"card-body\"><form method=\"POST\" action=\"/jobsform\"><div class=\"mb-3\"><label for=\"title\" class=\"form-label\">Title</label> <input type=\"text\" class=\"form-control\" id=\"title\" name=\"title\" required></div><div class=\"mb-3\"><label for=\"description\" class=\"form-label\">Description</label> <textarea class=\"form-control\" id=\"description\" name=\"description\" rows=\"5\" required></textarea></div><div class=\"mb-3\"><label for=\"category\" class=\"form-label\">Category</label> <select class=\"form-select\" id=\"category\" name=\"category_id\" required><option value=\"\" disabled selected>Select a category</option> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, category := range categories {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "<option value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d", category.ID))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/client/createjobs.templ`, Line: 79, Col: 82}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(category.Name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/client/createjobs.templ`, Line: 79, Col: 100}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "</select></div><div class=\"mb-3\"><label class=\"form-label\">Budget Range ($)</label><div class=\"d-flex gap-2\"><input type=\"number\" class=\"form-control\" name=\"budget_min\" placeholder=\"Min\" required> <input type=\"number\" class=\"form-control\" name=\"budget_max\" placeholder=\"Max\"></div></div><div class=\"mb-3\"><label for=\"deadline\" class=\"form-label\">Deadline</label> <input type=\"date\" class=\"form-control\" id=\"deadline\" name=\"deadline\"></div><button type=\"submit\" class=\"btn btn-dark w-100\">Post Job</button></form></div></div></div></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
