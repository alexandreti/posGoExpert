package utils

func ConverterTemperaturas(tempC float64) (float64, float64) {
	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15
	return tempF, tempK
}
