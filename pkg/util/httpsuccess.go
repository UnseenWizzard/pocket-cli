package util

func IsHttpSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}