package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

func PropertiesToJson(propertiesMap map[string]string) ([]byte, error) {
	resultJson := make(map[string]interface{})

	for key := range propertiesMap {
		createJsonObject(resultJson, key, propertiesMap[key])
	}

	return json.Marshal(resultJson)
}

func createJsonObject(resultJson map[string]interface{}, propertyKey, propertyValue string) map[string]interface{} {

	if !strings.Contains(propertyKey, ".") {
		resultJson[propertyKey] = propertyValue
		return resultJson
	}

	currentKey := getFirstKey(propertyKey)
	if len(currentKey) != 0 {

		println(fmt.Sprintf("\ncurrentKey = %s, propertyKey = %s", currentKey, propertyKey))
		runes := []rune(propertyKey)
		subRightKey := string(runes[len(currentKey)+1 : len(propertyKey)])
		fmt.Printf("\n subRightKey = %v", subRightKey)

		childJson := getJsonIfExists(resultJson, currentKey)
		resultJson[currentKey] = createJsonObject(childJson, subRightKey, propertyValue)
	}

	return resultJson
}

func getJsonIfExists(parent map[string]interface{}, key string) map[string]interface{} {
	if parent == nil {
		return nil
	}

	if parent[key] != nil {
		return parent[key].(map[string]interface{})
	} else {
		return make(map[string]interface{})
	}
}

func getFirstKey(fullKey string) string {
	splittedKey := strings.Split(fullKey, ".")

	fmt.Printf("\n fullKey = %v", fullKey)
	fmt.Printf("\n splittedKey = %v", splittedKey)
	if len(splittedKey) != 0 {
		return splittedKey[0]
	}

	return fullKey
}
