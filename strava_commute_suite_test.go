package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStravaCommute(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StravaCommute Suite")
}
