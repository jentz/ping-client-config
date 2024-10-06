package adminapi

import pfclient "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"

func GetSingleExtendedParameterValue(parameters *map[string]pfclient.ParameterValues, key string, defaultValue string) string {
	if parameters == nil {
		return defaultValue
	}
	parameterValues := (*parameters)[key]
	values, ok := parameterValues.GetValuesOk()
	if ok {
		return values[0]
	}
	return defaultValue
}
