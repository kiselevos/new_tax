package data

type Region struct {
	Label string
	Name  string
}

var regionLabelMap = map[string]string{
	// федеральные города
	"Moskva":          "msk",
	"Sankt-Peterburg": "spb",

	// республики
	"Adygeya, Respublika":                 "adygea",
	"Altay, Respublika":                   "altai_rep",
	"Bashkortostan, Respublika":           "bashkortostan",
	"Buryatiya, Respublika":               "buryatia",
	"Chechenskaya Respublika":             "chechnya",
	"Chuvashskaya Respublika":             "chuvashia",
	"Dagestan, Respublika":                "dagestan",
	"Ingushetiya, Respublika":             "ingushetia",
	"Kabardino-Balkarskaya Respublika":    "kabardino",
	"Kalmykiya, Respublika":               "kalmykia",
	"Karachayevo-Cherkesskaya Respublika": "karachay",
	"Kareliya, Respublika":                "karelia",
	"Khakasiya, Respublika":               "khakasiya",
	"Komi, Respublika":                    "komi",
	"Mariy El, Respublika":                "mari_el",
	"Mordoviya, Respublika":               "mordovia",
	"Saha, Respublika":                    "yakutia",
	"Severnaya Osetiya, Respublika":       "ossetia",
	"Tatarstan, Respublika":               "tatarstan",
	"Tyva, Respublika":                    "tuva",
	"Udmurtskaya Respublika":              "udmurtia",

	// края
	"Altayskiy kray":      "altai",
	"Kamchatskiy kray":    "kamchatka",
	"Khabarovskiy kray":   "khabarovsk",
	"Krasnodarskiy kray":  "krasnodar",
	"Krasnoyarskiy kray":  "krasnoyarsk",
	"Permskiy kray":       "perm",
	"Primorskiy kray":     "primorye",
	"Stavropol'skiy kray": "stavropol",
	"Zabaykal'skiy kray":  "zabaykalsky",

	// области
	"Amurskaya oblast'":        "amur",
	"Arkhangel'skaya oblast'":  "arkhangelsk",
	"Astrakhanskaya oblast'":   "astrakhan",
	"Belgorodskaya oblast'":    "belgorod",
	"Bryanskaya oblast'":       "bryansk",
	"Chelyabinskaya oblast'":   "chelyabinsk",
	"Irkutskaya oblast'":       "irkutsk",
	"Ivanovskaya oblast'":      "ivanovo",
	"Kaliningradskaya oblast'": "kaliningrad",
	"Kaluzhskaya oblast'":      "kaluga",
	"Kemerovskaya oblast'":     "kemerovo",
	"Kirovskaya oblast'":       "kirov",
	"Kostromskaya oblast'":     "kostroma",
	"Kurganskaya oblast'":      "kurgan",
	"Kurskaya oblast'":         "kursk",
	"Leningradskaya oblast'":   "leningrad",
	"Lipetskaya oblast'":       "lipetsk",
	"Magadanskaya oblast'":     "magadan",
	"Moskovskaya oblast'":      "moscow_obl",
	"Murmanskaya oblast'":      "murmansk",
	"Nizhegorodskaya oblast'":  "nizhny",
	"Novgorodskaya oblast'":    "novgorod",
	"Novosibirskaya oblast'":   "novosibirsk",
	"Omskaya oblast'":          "omsk",
	"Orenburgskaya oblast'":    "orenburg",
	"Orlovskaya oblast'":       "orel",
	"Penzenskaya oblast'":      "penza",
	"Pskovskaya oblast'":       "pskov",
	"Rostovskaya oblast'":      "rostov",
	"Ryazanskaya oblast'":      "ryazan",
	"Sakhalinskaya oblast'":    "sakhalin",
	"Samarskaya oblast'":       "samara",
	"Saratovskaya oblast'":     "saratov",
	"Smolenskaya oblast'":      "smolensk",
	"Sverdlovskaya oblast'":    "ekaterinburg",
	"Tambovskaya oblast'":      "tambov",
	"Tomskaya oblast'":         "tomsk",
	"Tul'skaya oblast'":        "tula",
	"Tverskaya oblast'":        "tver",
	"Tyumenskaya oblast'":      "tyumen",
	"Ul'yanovskaya oblast'":    "ulyanovsk",
	"Vladimirskaya oblast'":    "vladimir",
	"Volgogradskaya oblast'":   "volgograd",
	"Vologodskaya oblast'":     "vologda",
	"Voronezhskaya oblast'":    "voronezh",
	"Yaroslavskaya oblast'":    "yaroslavl",

	// автономные округа
	"Khanty-Mansiyskiy avtonomnyy okrug": "hmao",
	"Yamalo-Nenetskiy avtonomnyy okrug":  "yanao",
	"Nenetskiy avtonomnyy okrug":         "nao",
	"Chukotskiy avtonomnyy okrug":        "chukotka",
	"Yevreyskaya avtonomnaya oblast'":    "yeao",
}

func NormalizeRegion(name string) Region {
	label, ok := regionLabelMap[name]
	if !ok {
		label = "other"
	}

	return Region{
		Label: label,
		Name:  name,
	}
}
