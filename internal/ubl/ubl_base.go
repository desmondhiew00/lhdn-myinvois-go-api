package ubl

import (
	"github.com/desmondhiew00/lhdn-myinvois-go-api/internal/constant"
	"github.com/desmondhiew00/lhdn-myinvois-go-api/internal/util"
)

/* ------------------------------- Address UBL ------------------------------ */

type AddressUBLInput struct {
	CityName     string `json:"cityName"`
	PostalZone   string `json:"postalZone"`
	State        string `json:"state"`
	AddressLine0 string `json:"addressLine0"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	Country      string `json:"country"`
}

func AddressUBL(input AddressUBLInput) map[string]interface{} {
	return map[string]interface{}{
		"CityName": []interface{}{
			map[string]interface{}{"_": input.CityName},
		},
		"PostalZone": []interface{}{
			map[string]interface{}{"_": input.PostalZone},
		},
		"CountrySubentityCode": []interface{}{
			map[string]interface{}{"_": constant.Value(input.State, constant.StateCodes, "17")},
		},
		"AddressLine": []interface{}{
			map[string]interface{}{"Line": []interface{}{map[string]interface{}{"_": input.AddressLine0}}},
			map[string]interface{}{"Line": []interface{}{map[string]interface{}{"_": input.AddressLine1}}},
			map[string]interface{}{"Line": []interface{}{map[string]interface{}{"_": input.AddressLine2}}},
		},
		"Country": []interface{}{
			map[string]interface{}{
				"IdentificationCode": []interface{}{
					map[string]interface{}{
						"_":            constant.Value(input.Country, constant.CountryCodes, "MYS"),
						"listID":       "ISO3166-1",
						"listAgencyID": "6",
					},
				},
			},
		},
	}
}

/* ------------------------------- Contact UBL ------------------------------ */

func ContactUBL(contactNo, email string) map[string]interface{} {
	contact := make(map[string]interface{})
	if contactNo != "" {
		contact["Telephone"] = []interface{}{map[string]interface{}{"_": contactNo}}
	}
	if email != "" {
		contact["ElectronicMail"] = []interface{}{map[string]interface{}{"_": email}}
	}
	return contact
}

/* ------------------------------ Supplier UBL ------------------------------ */

type SupplierUBLInput struct {
	Name            string          `json:"name"`
	TIN             string          `json:"tin"`
	IdType          string          `json:"idType"`
	IdValue         string          `json:"idValue"` // NRIC, BRN, PASSPORT, ARMY
	SST             string          `json:"sst"`
	TTX             string          `json:"ttx"`
	ContactNumber   string          `json:"contactNumber"`
	Email           string          `json:"email"`
	MSICCode        string          `json:"msicCode"`
	MSICDescription string          `json:"msicDescription"`
	Address         AddressUBLInput `json:"address"`
}

func SupplierUBL(supplier SupplierUBLInput) map[string]interface{} {
	party := map[string]interface{}{
		"IndustryClassificationCode": []interface{}{
			map[string]interface{}{
				"_":            supplier.MSICCode,
				"name":         constant.Value(supplier.MSICDescription, constant.MSICDescription, "NA"),
				"listID":       "MSIC",
				"listAgencyID": "1",
			},
		},
		"PostalAddress": []interface{}{AddressUBL(supplier.Address)},
		"PartyLegalEntity": []interface{}{
			map[string]interface{}{
				"RegistrationName": []interface{}{
					map[string]interface{}{"_": supplier.Name},
				},
			},
		},
		"PartyIdentification": []interface{}{
			map[string]interface{}{
				"ID": []interface{}{map[string]interface{}{
					"_": util.NonEmpty(supplier.TIN, "NA"), "schemeID": "TIN",
				}},
			},
			map[string]interface{}{
				"ID": []interface{}{map[string]interface{}{
					"_": util.NonEmpty(supplier.IdValue, "NA"), "schemeID": util.NonEmpty(supplier.IdType, "BRN"),
				},
				},
			},
			map[string]interface{}{
				"ID": []interface{}{map[string]interface{}{
					"_": util.NonEmpty(supplier.SST, "NA"), "schemeID": "SST",
				},
				},
			},
			map[string]interface{}{
				"ID": []interface{}{map[string]interface{}{
					"_": util.NonEmpty(supplier.TTX, "NA"), "schemeID": "TTX",
				},
				},
			},
		},
		"Contact": []interface{}{ContactUBL(supplier.ContactNumber, supplier.Email)},
	}

	return map[string]interface{}{
		"AccountingSupplierParty": []interface{}{map[string]interface{}{
			"Party": []interface{}{party},
		}},
	}
}

/* ------------------------------ Buyer UBL ------------------------------ */

type BuyerUBLInput struct {
	Name                  string          `json:"name"`
	TIN                   string          `json:"tin"`
	IdType                string          `json:"idType"`
	IdValue               string          `json:"idValue"` // NRIC, BRN, PASSPORT, ARMY
	SstRegistrationNumber string          `json:"sstRegistrationNumber"`
	ContactNumber         string          `json:"contactNumber"`
	Email                 string          `json:"email"`
	Address               AddressUBLInput `json:"address"`
}

func BuyerUBL(buyer BuyerUBLInput) map[string]interface{} {
	party := map[string]interface{}{
		"PostalAddress": []interface{}{AddressUBL(buyer.Address)},
		"PartyLegalEntity": []interface{}{
			map[string]interface{}{
				"RegistrationName": []interface{}{map[string]interface{}{"_": buyer.Name}},
			},
		},
		"PartyIdentification": []interface{}{
			map[string]interface{}{"ID": []interface{}{map[string]interface{}{"_": util.NonEmpty(buyer.TIN, "NA"), "schemeID": "TIN"}}},
			map[string]interface{}{"ID": []interface{}{map[string]interface{}{"_": util.NonEmpty(buyer.IdValue, "NA"), "schemeID": util.NonEmpty(buyer.IdType, "NA")}}},
			map[string]interface{}{"ID": []interface{}{map[string]interface{}{"_": util.NonEmpty(buyer.SstRegistrationNumber, "NA"), "schemeID": "SST"}}},
			map[string]interface{}{"ID": []interface{}{map[string]interface{}{"_": "NA", "schemeID": "TTX"}}},
		},
		"Contact": []interface{}{ContactUBL(buyer.ContactNumber, buyer.Email)},
	}

	return map[string]interface{}{
		"AccountingCustomerParty": []interface{}{map[string]interface{}{
			"Party": []interface{}{party},
		}},
	}
}

func AddSupplierUBL(rawUbl map[string]interface{}, supplier SupplierUBLInput) {
	rawUbl["AccountingSupplierParty"] = SupplierUBL(supplier)["AccountingSupplierParty"]
}

func AddBuyerUBL(rawUbl map[string]interface{}, buyer BuyerUBLInput) {
	rawUbl["AccountingCustomerParty"] = BuyerUBL(buyer)["AccountingCustomerParty"]
}
