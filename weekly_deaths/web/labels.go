package web

var ageLabels = []MetadataLabel{
	{Value: "TOTAL", Label: "Total", Order: 1, Type: "age"},
	{Value: "UNK", Label: "Unknown", Order: 2, Type: "age"},
	{Value: "Y_LT5", Label: "<5", Order: 3, Type: "age"},
	{Value: "Y5-9", Label: "From 5 to 9", Order: 4, Type: "age"},
	{Value: "Y10-14", Label: "From 10 to 14", Order: 5, Type: "age"},
	{Value: "Y15-19", Label: "From 15 to 19", Order: 6, Type: "age"},
	{Value: "Y20-24", Label: "From 20 to 24", Order: 7, Type: "age"},
	{Value: "Y25-29", Label: "From 25 to 29", Order: 8, Type: "age"},
	{Value: "Y30-34", Label: "From 30 to 34", Order: 9, Type: "age"},
	{Value: "Y35-39", Label: "From 35 to 39", Order: 10, Type: "age"},
	{Value: "Y40-44", Label: "From 40 to 44", Order: 11, Type: "age"},
	{Value: "Y45-49", Label: "From 45 to 49", Order: 12, Type: "age"},
	{Value: "Y50-54", Label: "From 50 to 54", Order: 13, Type: "age"},
	{Value: "Y55-59", Label: "From 55 to 59", Order: 14, Type: "age"},
	{Value: "Y60-64", Label: "From 60 to 64", Order: 15, Type: "age"},
	{Value: "Y65-69", Label: "From 65 to 69", Order: 16, Type: "age"},
	{Value: "Y70-74", Label: "From 70 to 74", Order: 17, Type: "age"},
	{Value: "Y75-79", Label: "From 75 to 79", Order: 18, Type: "age"},
	{Value: "Y80-84", Label: "From 80 to 84", Order: 19, Type: "age"},
	{Value: "Y85-89", Label: "From 85 to 89", Order: 20, Type: "age"},
	{Value: "Y_GE90", Label: ">=90", Order: 21, Type: "age"},
}

var genderLabels = []MetadataLabel{
	{Value: "T", Label: "Total", Order: 1, Type: "gender"},
	{Value: "F", Label: "Female", Order: 2, Type: "gender"},
	{Value: "M", Label: "Male", Order: 3, Type: "gender"},
}

var countryLabels = []MetadataLabel{
	{Value: "AD", Label: "Andorra", Order: 1, Type: "country"},
	{Value: "AL", Label: "Albania", Order: 2, Type: "country"},
	{Value: "AM", Label: "Armenia", Order: 3, Type: "country"},
	{Value: "AT", Label: "Austria", Order: 4, Type: "country"},
	{Value: "BE", Label: "Belgium", Order: 5, Type: "country"},
	{Value: "BG", Label: "Bulgaria", Order: 6, Type: "country"},
	{Value: "CH", Label: "Switzerland", Order: 7, Type: "country"},
	{Value: "CY", Label: "Cyprus", Order: 8, Type: "country"},
	{Value: "CZ", Label: "Czechia", Order: 9, Type: "country"},
	{Value: "DE", Label: "Germany", Order: 10, Type: "country"},
	{Value: "DK", Label: "Denmark", Order: 11, Type: "country"},
	{Value: "EE", Label: "Estonia", Order: 12, Type: "country"},
	{Value: "EL", Label: "Greece", Order: 13, Type: "country"},
	{Value: "ES", Label: "Spain", Order: 14, Type: "country"},
	{Value: "FI", Label: "Finland", Order: 15, Type: "country"},
	{Value: "FR", Label: "France", Order: 16, Type: "country"},
	{Value: "GE", Label: "Georgia", Order: 17, Type: "country"},
	{Value: "HR", Label: "Croatia", Order: 18, Type: "country"},
	{Value: "HU", Label: "Hungary", Order: 19, Type: "country"},
	{Value: "IE", Label: "Ireland", Order: 20, Type: "country"},
	{Value: "IS", Label: "Iceland", Order: 21, Type: "country"},
	{Value: "IT", Label: "Italy", Order: 22, Type: "country"},
	{Value: "LI", Label: "Liechtenstein", Order: 23, Type: "country"},
	{Value: "LT", Label: "Lithuania", Order: 24, Type: "country"},
	{Value: "LU", Label: "Luxembourg", Order: 25, Type: "country"},
	{Value: "LV", Label: "Latvia", Order: 26, Type: "country"},
	{Value: "ME", Label: "Montenegro", Order: 27, Type: "country"},
	{Value: "MT", Label: "Malta", Order: 28, Type: "country"},
	{Value: "NL", Label: "Netherlands", Order: 29, Type: "country"},
	{Value: "NO", Label: "Norway", Order: 30, Type: "country"},
	{Value: "PL", Label: "Poland", Order: 31, Type: "country"},
	{Value: "PT", Label: "Portugal", Order: 32, Type: "country"},
	{Value: "RO", Label: "Romania", Order: 33, Type: "country"},
	{Value: "RS", Label: "Serbia", Order: 34, Type: "country"},
	{Value: "SE", Label: "Sweden", Order: 35, Type: "country"},
	{Value: "SI", Label: "Slovenia", Order: 36, Type: "country"},
	{Value: "SK", Label: "Slovakia", Order: 37, Type: "country"},
	{Value: "UK", Label: "United Kingdom", Order: 38, Type: "country"},
}

// GetLabels returns all static labels (country, age, gender)
// for data contained within Eurostat Weekly Deaths dataset.
func GetLabels() []MetadataLabel {
	labels := make([]MetadataLabel, 0)
	data := [][]MetadataLabel{ageLabels, countryLabels, genderLabels}

	for _, d := range data {
		labels = append(labels, d...)
	}

	return labels
}
