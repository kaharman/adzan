package domain

import "fmt"

var DayEn = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
var DayId = []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
var DayAr = []string{"الأحد", "الاثنين", "الثلاثاء", "الأربعاء", "الخميس", "الجمعة", "السبت"}
var HijriMonthEn = []string{"al-Muḥarram", "Ṣafar", "Rabīʿ al-ʾAwwal", "Rabīʿ ath-Thānī", "Jumadā al-ʾŪlā", "Jumādā ath-Thāniyah", "Rajab", "Shaʿbān", "Ramaḍān", "Shawwāl", "Ḏū al-Qaʿdah", "Ḏū al-Ḥijjah"}
var HijriMonthId = []string{"Muḥarram", "Ṣafar", "Rabīʿul 'Awwal", "Rabīʿuṡ Ṡānī", "Jumadal 'Ūlā", "Jumādaṡ Ṡāniyah", "Rajab", "Syaʿbān", "Ramaḍān", "Syawwāl", "Żulqaʿdah", "Żulḥijjah"}
var HijriMonthAr = []string{"ٱلْمُحَرَّم", "صَفَر", "رَبِيع ٱلْأَوَّل", "رَبِيع ٱلثَّانِي", "جُمَادَىٰ ٱلْأُولَىٰ", "جُمَادَىٰ ٱلثَّانِيَة", "رَجَب", "شَعْبَان", "رَمَضَان", "شَوَّال", "ذُو ٱلْقَعْدَة", "ذُو ٱلْحِجَّة"}

func ConvertDigitEnglishToArabic(number string) string {
	var result string

	for _, char := range number {
		result = fmt.Sprintf("%s%c", result, rune(char+0x630))
	}

	return result
}
