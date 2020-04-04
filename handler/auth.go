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

			if len(username) < 3 || !isTokenVaild(token) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			// 调用传入的h来继续服务该请求
			h(w, r)
		})
}
