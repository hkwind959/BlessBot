package common

import (
	"math/rand"
	"time"
)

// GenerateRandomHardwareInfo 生成随机的硬件信息
func GenerateRandomHardwareInfo() map[string]interface{} {
	cpuArchitectures := []string{"x86_64", "ARM64", "x86"}
	cpuModels := []string{"Intel Core i7-10700K CPU @ 3.80GHz",
		"AMD Ryzen 5 5600G with Radeon Graphics",
		"Intel Core i5-10600K CPU @ 4.10GHz",
		"AMD Ryzen 7 5800X",
		"Intel Core i9-10900K CPU @ 3.70GHz",
		"AMD Ryzen 9 5900X",
		"Intel Core i3-10100 CPU @ 3.60GHz",
		"AMD Ryzen 3 3300X",
		"Intel Core i7-9700K CPU @ 3.60GHz",
		"12th Gen Intel(R) Core(TM) i7-12700H",
	}

	cpuFeatures := []string{"mmx", "sse", "sse2", "sse3", "ssse3", "sse4_1", "sse4_2", "avx", "avx2", "fma"}
	// CPU核数
	numProcessors := []int{4, 6, 8, 12, 16}
	// 内存大小
	memorySizes := []int{8 * 1024 * 1024 * 1024, 16 * 1024 * 1024 * 1024, 32 * 1024 * 1024 * 1024, 64 * 1024 * 1024 * 1024}
	randomCpuFeatures := make([]string, 0)
	featureCount := rand.Intn(len(cpuFeatures)) + 1
	for i := 0; i < featureCount; i++ {
		randomCpuFeatures = append(randomCpuFeatures, getRandomElement(cpuFeatures))
	}

	uniqueFeatures := make(map[string]bool)
	for _, feature := range randomCpuFeatures {
		uniqueFeatures[feature] = true
	}

	var uniqueFeatureList []string
	for feature := range uniqueFeatures {
		uniqueFeatureList = append(uniqueFeatureList, feature)
	}

	return map[string]interface{}{
		"cpuArchitecture": getRandomElement(cpuArchitectures),
		"cpuModel":        getRandomElement(cpuModels),
		"cpuFeatures":     uniqueFeatureList,
		"numOfProcessors": getRandomElementInt(numProcessors),
		"totalMemory":     getRandomElementInt(memorySizes),
	}
}

func getRandomElement(slice []string) string {
	rand.Seed(time.Now().UnixNano())
	return slice[rand.Intn(len(slice))]
}

func getRandomElementInt(slice []int) int {
	rand.Seed(time.Now().UnixNano())
	return slice[rand.Intn(len(slice))]
}
