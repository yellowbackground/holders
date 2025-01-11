package examples

import "github.com/yellowbackground/holders"

var Collections = []holders.Collection{
	{
		Name: "Yieldlings Flambos",
		Addresses: []string{
			"5DYIZMX7N4SAB44HLVRUGLYBPSN4UMPDZVTX7V73AIRMJQA3LKTENTLFZ4",
		},
		IncludeNameContains: []string{"Flamborghini"},
	},
	{
		Name: "Yieldlings",
		Addresses: []string{
			"5DYIZMX7N4SAB44HLVRUGLYBPSN4UMPDZVTX7V73AIRMJQA3LKTENTLFZ4",
		},
		ExcludeNameContains: []string{"Flamborghini"},
		UnitNamePrefixes:    []string{"TLDG", "YLD"},
	},
	{
		Name: "M.N.G.O",
		Addresses: []string{
			"MNGOLDXO723TDRM6527G7OZ2N7JLNGCIH6U2R4MOCPPLONE3ZATOBN7OQM",
			"MNGORTG4A3SLQXVRICQXOSGQ7CPXUPMHZT3FJZBIZHRYAQCYMEW6VORBIA",
			"MNGOZ3JAS3C4QTGDQ5NVABUEZIIF4GAZY52L3EZE7BQIBFTZCNLQPXHRHE",
			"MNGO4JTLBN64PJLWTQZYHDMF2UBHGJGW5L7TXDVTJV7JGVD5AE4Y3HTEZM",
		},
		UnitNamePrefixes: []string{"MNGO"},
	},
	{
		Name: "Mostly Frens",
		Addresses: []string{
			"MOSTLYSNUJP7PG6Q3FNJCGGENQXMOH3PXXMIJRFLODLG2DNDBHI7QHJSOE",
		},
		UnitNamePrefixes: []string{"MFER"},
	},
	{
		Name: "Best Frens",
		Addresses: []string{
			"MOSTLYSNUJP7PG6Q3FNJCGGENQXMOH3PXXMIJRFLODLG2DNDBHI7QHJSOE",
		},
		UnitNamePrefixes: []string{"BFER"},
	},
}
