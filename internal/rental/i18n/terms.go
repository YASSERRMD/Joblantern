// Package i18n holds rental-specific terminology lookups across the
// eight supported languages. Keys match the IDs used in the rental UI
// templates and rule pack.
package i18n

// Terms maps a language tag (BCP-47) to a label table.
var Terms = map[string]map[string]string{
	"en": {
		"deposit": "Security deposit", "landlord": "Landlord", "lease": "Lease", "viewing": "Viewing",
	},
	"hi": {
		"deposit": "जमानत राशि", "landlord": "मकान मालिक", "lease": "किरायानामा", "viewing": "देखने का समय",
	},
	"bn": {
		"deposit": "সিকিউরিটি ডিপোজিট", "landlord": "মালিক", "lease": "ভাড়াপত্র", "viewing": "পরিদর্শন",
	},
	"tl": {
		"deposit": "Deposito", "landlord": "May-ari", "lease": "Kontrata sa upa", "viewing": "Pagbisita",
	},
	"id": {
		"deposit": "Uang jaminan", "landlord": "Pemilik", "lease": "Sewa", "viewing": "Kunjungan",
	},
	"ar": {
		"deposit": "تأمين", "landlord": "المالك", "lease": "عقد إيجار", "viewing": "معاينة",
	},
	"ur": {
		"deposit": "ضمانتی رقم", "landlord": "مالک مکان", "lease": "کرایہ نامہ", "viewing": "معائنہ",
	},
	"ne": {
		"deposit": "धरौटी", "landlord": "घरबेटी", "lease": "बहाल सम्झौता", "viewing": "अवलोकन",
	},
}
