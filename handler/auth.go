package handler

import "net/http"

// HTTPInterceptor : token 拦截器，验证username和token是否有效
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()

			// 获取用户名和密码
			username := r.Form.Get("username")
			token := r.Form.Get("token")

			if len(username) < 3 || !isTokenVaild(username, token) {
				http.Redirect(w, r, "/user/signin", http.StatusFound)
				// w.WriteHeader(http.StatusForbidden)
				return
			}
			// 调用传入的h来继续服务该请求
			h(w, r)
		})
}

// AgreeInterceptor : 同意拦截器，用来拦截当前的用户是否可以具有权限同意
func AgreeInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()

			// 获取用户名和密码
			username := r.Form.Get("username")
			token := r.Form.Get("token")

			if len(username) < 3 || !isTokenVaild(username, token) {
				http.Redirect(w, r, "/user/signin", http.StatusFound)
				// w.WriteHeader(http.StatusForbidden)
				return
			}
			// 调用传入的h来继续服务该请求
			h(w, r)
		})
}
