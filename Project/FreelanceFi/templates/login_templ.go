// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

// login.templ
func LoginPage() templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><title>HTMX + Go Login</title><link href=\"https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css\" rel=\"stylesheet\"><script src=\"https://unpkg.com/htmx.org@1.9.2\"></script></head><body class=\"d-flex justify-content-center align-items-center vh-100\"><div class=\"card p-4 shadow-sm\" style=\"width:350px;\"><h3 class=\"mb-3 text-center\">Login</h3><form id=\"login-form\" hx-post=\"/login\" hx-target=\"#response\" hx-swap=\"innerHTML\" class=\"needs-validation\" novalidate><div id=\"response\" class=\"mb-3\"></div><div class=\"mb-3\"><label class=\"form-label\">Username</label> <input type=\"text\" name=\"username\" class=\"form-control\" required><div class=\"invalid-feedback\">Please enter a valid username.</div></div><div class=\"mb-3\"><label class=\"form-label\">Password</label> <input type=\"password\" name=\"password\" class=\"form-control\" required><div class=\"invalid-feedback\">Please enter your password.</div></div><button class=\"btn btn-primary w-100\">Login</button></form><div class=\"mt-3 text-center\"><a href=\"/register\" class=\"link-primary\">Don't have an account? Register</a></div></div><script>\n        (function(){\n            'use strict';\n            let form = document.getElementById('login-form');\n\n            form.addEventListener('submit', function(e) {\n                if (!form.checkValidity()) {\n                    e.preventDefault(); e.stopPropagation();\n                }\n                form.classList.add('was-validated');\n            }, false);\n\n            form.addEventListener('htmx:beforeSend', function(evt) {\n                const inputs = form.querySelectorAll('input');\n                inputs.forEach(input => {\n                    input.classList.remove('is-valid', 'is-invalid');\n                    if (input.nextElementSibling) {\n                        input.nextElementSibling.textContent = '';\n                    }\n                });\n            });\n\n            form.addEventListener('htmx:afterSwap', function(evt) {\n                if (evt.detail.target.id === 'response') {\n                    const usernameInput = form.querySelector('[name=username]');\n                    const passwordInput = form.querySelector('[name=password]');\n\n                    if (evt.detail.xhr.status === 401 || evt.detail.xhr.status === 400) {\n                        usernameInput.classList.add('is-invalid');\n                        passwordInput.classList.add('is-invalid');\n                        if (usernameInput.nextElementSibling) {\n                            usernameInput.nextElementSibling.textContent = \"Incorrect username.\";\n                        }\n                        if (passwordInput.nextElementSibling) {\n                            passwordInput.nextElementSibling.textContent = \"Incorrect password.\";\n                        }\n                    }\n\n                    const alert = document.querySelector('#response .alert');\n                    if (alert) {\n                        setTimeout(() => {\n                            alert.classList.add('fade');\n                            alert.classList.add('show');\n                            setTimeout(() => {\n                                alert.remove();\n                            }, 1000);\n                        }, 3000);\n                    }\n                }\n            });\n        })();\n    </script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
