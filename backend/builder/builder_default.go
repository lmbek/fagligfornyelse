//go:build !development

package builder

func IsDevelopmentMode() bool {
	return false
}
