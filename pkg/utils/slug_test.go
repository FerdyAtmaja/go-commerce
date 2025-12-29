package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple text",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "Text with special characters",
			input:    "Hello, World! & More",
			expected: "hello-world-more",
		},
		{
			name:     "Text with numbers",
			input:    "iPhone 15 Pro Max",
			expected: "iphone-15-pro-max",
		},
		{
			name:     "Text with multiple spaces",
			input:    "Multiple   Spaces   Here",
			expected: "multiple-spaces-here",
		},
		{
			name:     "Text with leading/trailing spaces",
			input:    "  Trimmed Text  ",
			expected: "trimmed-text",
		},
		{
			name:     "Indonesian text",
			input:    "Kategori Elektronik & Gadget",
			expected: "kategori-elektronik-gadget",
		},
		{
			name:     "Text with underscores",
			input:    "snake_case_text",
			expected: "snake-case-text",
		},
		{
			name:     "Already lowercase",
			input:    "already lowercase",
			expected: "already-lowercase",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only special characters",
			input:    "!@#$%^&*()",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSlug(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnsureUniqueSlug(t *testing.T) {
	t.Run("Slug is unique", func(t *testing.T) {
		baseSlug := "unique-slug"
		checkExists := func(slug string) bool {
			return false // Always return false, meaning slug doesn't exist
		}

		result := EnsureUniqueSlug(baseSlug, checkExists)
		assert.Equal(t, baseSlug, result)
	})

	t.Run("Slug exists once", func(t *testing.T) {
		baseSlug := "existing-slug"
		checkExists := func(slug string) bool {
			return slug == "existing-slug" // Only base slug exists
		}

		result := EnsureUniqueSlug(baseSlug, checkExists)
		assert.Equal(t, "existing-slug-1", result)
	})

	t.Run("Slug exists multiple times", func(t *testing.T) {
		baseSlug := "popular-slug"
		existingSlugs := map[string]bool{
			"popular-slug":   true,
			"popular-slug-1": true,
			"popular-slug-2": true,
		}
		
		checkExists := func(slug string) bool {
			return existingSlugs[slug]
		}

		result := EnsureUniqueSlug(baseSlug, checkExists)
		assert.Equal(t, "popular-slug-3", result)
	})

	t.Run("Empty base slug", func(t *testing.T) {
		baseSlug := ""
		checkExists := func(slug string) bool {
			return slug == "" // Empty slug exists
		}

		result := EnsureUniqueSlug(baseSlug, checkExists)
		assert.Equal(t, "-1", result)
	})
}

func TestGenerateSlugIntegration(t *testing.T) {
	// Test realistic scenarios
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Product name",
			input:    "Samsung Galaxy S24 Ultra 256GB",
			expected: "samsung-galaxy-s24-ultra-256gb",
		},
		{
			name:     "Category name",
			input:    "Komputer & Laptop",
			expected: "komputer-laptop",
		},
		{
			name:     "Store name",
			input:    "Toko Elektronik Jaya Abadi",
			expected: "toko-elektronik-jaya-abadi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSlug(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}