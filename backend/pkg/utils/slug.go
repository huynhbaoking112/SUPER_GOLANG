package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func GenerateSlug(name string) string {
	slug := strings.ToLower(name)

	slug = strings.ReplaceAll(slug, " ", "-")

	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	slug = reg.ReplaceAllString(slug, "")

	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	if slug == "" {
		slug = "workspace"
	}

	return slug
}

func GenerateUniqueSlug(name string, checkExists func(string) (bool, error)) (string, error) {
	baseSlug := GenerateSlug(name)

	exists, err := checkExists(baseSlug)
	if err != nil {
		return "", fmt.Errorf("failed to check slug existence: %w", err)
	}

	if !exists {
		return baseSlug, nil
	}

	for i := 1; i <= 1000; i++ {
		candidateSlug := fmt.Sprintf("%s-%d", baseSlug, i)

		exists, err := checkExists(candidateSlug)
		if err != nil {
			return "", fmt.Errorf("failed to check slug existence for %s: %w", candidateSlug, err)
		}

		if !exists {
			return candidateSlug, nil
		}
	}

	return "", fmt.Errorf("unable to generate unique slug for name: %s", name)
}
