package dto

type BanyanGetLoadDetailsDto struct {
	LoadDetails []struct {
		Quotes []struct {
			LoadID         int    `json:"LoadID"`
			QuoteID        int    `json:"QuoteID"`
			CarrierName    string `json:"CarrierName"`
			SCAC           string `json:"SCAC"`
			ThirdPartyName string `json:"ThirdPartyName"`
			ThirdPartySCAC string `json:"ThirdPartySCAC"`
			ServiceID      string `json:"ServiceID"`
			RawPrice       struct {
				NetPrice           float64 `json:"NetPrice"`
				FreightCharge      float64 `json:"FreightCharge"`
				DiscountAmount     int     `json:"DiscountAmount"`
				FuelSurcharge      float64 `json:"FuelSurcharge"`
				Minimum            int     `json:"Minimum"`
				GrossCharge        float64 `json:"GrossCharge"`
				Tariff             int     `json:"Tariff"`
				Interline          int     `json:"Interline"`
				AccessorialCharges float64 `json:"AccessorialCharges"`
				OtherCharges       int     `json:"OtherCharges"`
				Charges            []struct {
					Name   string  `json:"Name"`
					Amount float64 `json:"Amount"`
					Code   string  `json:"Code"`
				} `json:"Charges"`
				Markup int `json:"Markup"`
			} `json:"RawPrice"`
			CarrierPrice struct {
				NetPrice           float64 `json:"NetPrice"`
				FreightCharge      float64 `json:"FreightCharge"`
				DiscountAmount     int     `json:"DiscountAmount"`
				FuelSurcharge      float64 `json:"FuelSurcharge"`
				Minimum            int     `json:"Minimum"`
				GrossCharge        float64 `json:"GrossCharge"`
				Tariff             int     `json:"Tariff"`
				Interline          int     `json:"Interline"`
				AccessorialCharges float64 `json:"AccessorialCharges"`
				OtherCharges       float64 `json:"OtherCharges"`
				Charges            []struct {
					Name   string  `json:"Name"`
					Amount float64 `json:"Amount"`
					Code   string  `json:"Code"`
				} `json:"Charges"`
				Markup float64 `json:"Markup"`
			} `json:"CarrierPrice"`
			CustomerPrice struct {
				NetPrice           float64 `json:"NetPrice"`
				FreightCharge      float64 `json:"FreightCharge"`
				DiscountAmount     int     `json:"DiscountAmount"`
				FuelSurcharge      float64 `json:"FuelSurcharge"`
				Minimum            int     `json:"Minimum"`
				GrossCharge        float64 `json:"GrossCharge"`
				Tariff             int     `json:"Tariff"`
				Interline          int     `json:"Interline"`
				AccessorialCharges float64 `json:"AccessorialCharges"`
				OtherCharges       float64 `json:"OtherCharges"`
				Charges            []struct {
					Name   string  `json:"Name"`
					Amount float64 `json:"Amount"`
					Code   string  `json:"Code"`
				} `json:"Charges"`
				Markup float64 `json:"Markup"`
			} `json:"CustomerPrice"`
			TransitTime      int         `json:"TransitTime"`
			QuoteNumber      string      `json:"QuoteNumber"`
			CarrierPerson    string      `json:"CarrierPerson"`
			CarrierNote      string      `json:"CarrierNote"`
			Datestamp        string      `json:"Datestamp"`
			Interline        bool        `json:"Interline"`
			Accepted         bool        `json:"Accepted"`
			CurrencyType     int         `json:"CurrencyType"`
			Service          int         `json:"Service"`
			InsuranceDetails interface{} `json:"InsuranceDetails"`
			AccountNumber    string      `json:"AccountNumber"`
		} `json:"Quotes"`
		ClientName string `json:"ClientName"`
		Status     string `json:"Status"`
		Notes      []struct {
			Text      string `json:"Text"`
			User      string `json:"User"`
			DateStamp string `json:"DateStamp"`
			NoteType  int    `json:"NoteType"`
		} `json:"Notes"`
		Loadinfo struct {
			LoadID                         int         `json:"LoadID"`
			ManifestID                     string      `json:"ManifestID"`
			BOLNumber                      string      `json:"BOLNumber"`
			CustomerPO                     string      `json:"CustomerPO"`
			InvoiceID                      string      `json:"InvoiceID"`
			BillingID                      string      `json:"BillingID"`
			IncoTermID                     interface{} `json:"IncoTermID"`
			UltimateDestinationCountryCode string      `json:"UltimateDestinationCountryCode"`
			PickupNumber                   interface{} `json:"PickupNumber"`
			EstimatedDeliveryDate          string      `json:"EstimatedDeliveryDate"`
			EstimatedPickupDate            string      `json:"EstimatedPickupDate"`
			ActualPickupDate               string      `json:"ActualPickupDate"`
			ActualDeliveryDate             string      `json:"ActualDeliveryDate"`
		} `json:"Loadinfo"`
		BillTo struct {
			Name             string `json:"Name"`
			Note             string `json:"Note"`
			ShipType         int    `json:"ShipType"`
			PayType          int    `json:"PayType"`
			UseDefaultBillTo bool   `json:"UseDefaultBillTo"`
			AddressInfo      struct {
				Address1     string `json:"Address1"`
				Address2     string `json:"Address2"`
				City         string `json:"City"`
				CountryName  string `json:"CountryName"`
				CountryCode  string `json:"CountryCode"`
				State        string `json:"State"`
				Zipcode      string `json:"Zipcode"`
				LocationName string `json:"LocationName"`
			} `json:"AddressInfo"`
			ContactInfo struct {
				FirstName   string `json:"FirstName"`
				LastName    string `json:"LastName"`
				ContactName string `json:"ContactName"`
				Phone       string `json:"Phone"`
				PhoneExt    string `json:"PhoneExt"`
				Fax         string `json:"Fax"`
				Email       string `json:"Email"`
			} `json:"ContactInfo"`
			UseDefaultShipType bool `json:"UseDefaultShipType"`
			UseDefaultPayType  bool `json:"UseDefaultPayType"`
		} `json:"BillTo"`
		InsuranceInfo interface{} `json:"InsuranceInfo"`
		RateServices  []struct {
			ServiceCode         int    `json:"ServiceCode"`
			ShippingQty         int    `json:"ShippingQty"`
			PackageType         int    `json:"PackageType"`
			EquipmentType       int    `json:"EquipmentType"`
			AdditionalWeight    int    `json:"AdditionalWeight"`
			SpecialInstructions string `json:"SpecialInstructions"`
			Length              int    `json:"Length"`
			Width               int    `json:"Width"`
			Height              int    `json:"Height"`
			WeightUom           int    `json:"WeightUom"`
			SizeUom             int    `json:"SizeUom"`
		} `json:"RateServices"`
		Shipper struct {
			ContactInfo struct {
				FirstName   string `json:"FirstName"`
				LastName    string `json:"LastName"`
				ContactName string `json:"ContactName"`
				Phone       string `json:"Phone"`
				PhoneExt    string `json:"PhoneExt"`
				Fax         string `json:"Fax"`
				Email       string `json:"Email"`
			} `json:"ContactInfo"`
			AddressInfo struct {
				Address1     string `json:"Address1"`
				Address2     string `json:"Address2"`
				City         string `json:"City"`
				CountryName  string `json:"CountryName"`
				CountryCode  string `json:"CountryCode"`
				State        string `json:"State"`
				Zipcode      string `json:"Zipcode"`
				LocationName string `json:"LocationName"`
			} `json:"AddressInfo"`
			CompanyName  string      `json:"CompanyName"`
			Note         string      `json:"Note"`
			CompanyID    string      `json:"CompanyID"`
			LocationName interface{} `json:"LocationName"`
			VendorID     string      `json:"VendorID"`
			DCRefNum     interface{} `json:"DCRefNum"`
			Dock         struct {
				Name               string `json:"Name"`
				Note               string `json:"Note"`
				OpenTime           string `json:"OpenTime"`
				ShipmentDateTime   string `json:"ShipmentDateTime"`
				CloseTime          string `json:"CloseTime"`
				ConfirmationNumber string `json:"ConfirmationNumber"`
				FCFS               bool   `json:"FCFS"`
			} `json:"Dock"`
		} `json:"Shipper"`
		Consignee struct {
			ContactInfo struct {
				FirstName   string `json:"FirstName"`
				LastName    string `json:"LastName"`
				ContactName string `json:"ContactName"`
				Phone       string `json:"Phone"`
				PhoneExt    string `json:"PhoneExt"`
				Fax         string `json:"Fax"`
				Email       string `json:"Email"`
			} `json:"ContactInfo"`
			AddressInfo struct {
				Address1     string `json:"Address1"`
				Address2     string `json:"Address2"`
				City         string `json:"City"`
				CountryName  string `json:"CountryName"`
				CountryCode  string `json:"CountryCode"`
				State        string `json:"State"`
				Zipcode      string `json:"Zipcode"`
				LocationName string `json:"LocationName"`
			} `json:"AddressInfo"`
			CompanyName  string      `json:"CompanyName"`
			Note         string      `json:"Note"`
			CompanyID    string      `json:"CompanyID"`
			LocationName interface{} `json:"LocationName"`
			VendorID     string      `json:"VendorID"`
			DCRefNum     interface{} `json:"DCRefNum"`
			Dock         struct {
				Name               string `json:"Name"`
				Note               string `json:"Note"`
				OpenTime           string `json:"OpenTime"`
				ShipmentDateTime   string `json:"ShipmentDateTime"`
				CloseTime          string `json:"CloseTime"`
				ConfirmationNumber string `json:"ConfirmationNumber"`
				FCFS               bool   `json:"FCFS"`
			} `json:"Dock"`
		} `json:"Consignee"`
		ReturnLocation interface{} `json:"ReturnLocation"`
		PackageInfo    struct {
			CODAmount         int    `json:"CODAmount"`
			DeclaredLiability int    `json:"DeclaredLiability"`
			RouteNumber       string `json:"RouteNumber"`
		} `json:"PackageInfo"`
		Products []struct {
			Quantity          int    `json:"Quantity"`
			PackageType       int    `json:"PackageType"`
			Weight            int    `json:"Weight"`
			Class             int    `json:"Class"`
			NMFC              string `json:"NMFC"`
			SKU               string `json:"SKU"`
			IsHazmat          bool   `json:"IsHazmat"`
			HazmatPhoneNumber string `json:"HazmatPhoneNumber"`
			HazmatPhoneExt    string `json:"HazmatPhoneExt"`
			Description       string `json:"Description"`
			Length            int    `json:"Length"`
			Width             int    `json:"Width"`
			Height            int    `json:"Height"`
			UOM               int    `json:"UOM"`
			SortOrder         int    `json:"SortOrder"`
			ReferenceNumber   string `json:"ReferenceNumber"`
			ParcelOptions     struct {
				DeliveryConfirmation int         `json:"DeliveryConfirmation"`
				COD                  int         `json:"COD"`
				CODAmount            int         `json:"CODAmount"`
				AdditionalHandling   bool        `json:"AdditionalHandling"`
				LargePackage         bool        `json:"LargePackage"`
				DeclaredValue        int         `json:"DeclaredValue"`
				UnitValue            interface{} `json:"UnitValue"`
				UnitWeight           interface{} `json:"UnitWeight"`
				UnitType             interface{} `json:"UnitType"`
			} `json:"ParcelOptions"`
			UnNumber             string `json:"UnNumber"`
			HazMatShippingName   string `json:"HazMatShippingName"`
			HazMatPkgGroup       string `json:"HazMatPkgGroup"`
			HazMatClass          string `json:"HazMatClass"`
			InternationalOptions struct {
				ScheduleBCode         interface{} `json:"ScheduleBCode"`
				HarmonizedSystemCode  interface{} `json:"HarmonizedSystemCode"`
				ECCN                  interface{} `json:"ECCN"`
				UnitOriginCountryCode interface{} `json:"UnitOriginCountryCode"`
			} `json:"InternationalOptions"`
			HazMatContact string `json:"HazMatContact"`
			HazMatPoison  bool   `json:"HazMatPoison"`
		} `json:"Products"`
		ShipperAccessorials struct {
			AppointmentRequired   bool   `json:"AppointmentRequired"`
			InsidePickup          bool   `json:"InsidePickup"`
			SortSegregate         bool   `json:"SortSegregate"`
			PalletJack            bool   `json:"PalletJack"`
			ResidentialPickup     bool   `json:"ResidentialPickup"`
			LiftgatePickup        bool   `json:"LiftgatePickup"`
			MarkingTagging        bool   `json:"MarkingTagging"`
			TradeShowPickup       bool   `json:"TradeShowPickup"`
			NYCMetro              bool   `json:"NYCMetro"`
			NonBusinessHourPickup bool   `json:"NonBusinessHourPickup"`
			LimitedAccessType     string `json:"LimitedAccessType"`
		} `json:"ShipperAccessorials"`
		ConsigneeAccessorials struct {
			AppointmentRequired     bool   `json:"AppointmentRequired"`
			InsideDelivery          bool   `json:"InsideDelivery"`
			SortSegregate           bool   `json:"SortSegregate"`
			PalletJack              bool   `json:"PalletJack"`
			ResidentialDelivery     bool   `json:"ResidentialDelivery"`
			LiftgateDelivery        bool   `json:"LiftgateDelivery"`
			MarkingTagging          bool   `json:"MarkingTagging"`
			TradeShowDelivery       bool   `json:"TradeShowDelivery"`
			NYCMetro                bool   `json:"NYCMetro"`
			DeliveryNotification    bool   `json:"DeliveryNotification"`
			TwoHourSpecialDelivery  bool   `json:"TwoHourSpecialDelivery"`
			NonBusinessHourDelivery bool   `json:"NonBusinessHourDelivery"`
			LimitedAccessType       string `json:"LimitedAccessType"`
		} `json:"ConsigneeAccessorials"`
		LoadAccessorials struct {
			Guaranteed                   bool `json:"Guaranteed"`
			TimeDefinite                 bool `json:"TimeDefinite"`
			Expedited                    bool `json:"Expedited"`
			HolidayPickup                bool `json:"HolidayPickup"`
			HolidayDelivery              bool `json:"HolidayDelivery"`
			WeightDetermination          bool `json:"WeightDetermination"`
			BlindShipment                bool `json:"BlindShipment"`
			BlanketService               bool `json:"BlanketService"`
			ProtectFromFreezing          bool `json:"ProtectFromFreezing"`
			SingleShipment               bool `json:"SingleShipment"`
			CustomsInBond                bool `json:"CustomsInBond"`
			OverDimension                bool `json:"OverDimension"`
			Stackable                    bool `json:"Stackable"`
			Turnkey                      bool `json:"Turnkey"`
			FoodGradeProducts            bool `json:"FoodGradeProducts"`
			TSA                          bool `json:"TSA"`
			Bulkhead                     bool `json:"Bulkhead"`
			SignatureRequired            bool `json:"SignatureRequired"`
			BlanketServiceChilled        bool `json:"BlanketServiceChilled"`
			BlanketServiceFrozen         bool `json:"BlanketServiceFrozen"`
			SaturdayDelivery             bool `json:"SaturdayDelivery"`
			SecondMan                    bool `json:"SecondMan"`
			ReturnReceipt                bool `json:"ReturnReceipt"`
			ShipmentHold                 bool `json:"ShipmentHold"`
			ProactiveResponse            bool `json:"ProactiveResponse"`
			ShipperRelease               bool `json:"ShipperRelease"`
			WhiteGlove                   bool `json:"WhiteGlove"`
			RestrictedDelivery           bool `json:"RestrictedDelivery"`
			TankerEndorsedDriverRequired bool `json:"TankerEndorsedDriverRequired"`
		} `json:"LoadAccessorials"`
		UserDefined    []interface{} `json:"UserDefined"`
		ReferenceField []interface{} `json:"ReferenceField"`
	} `json:"LoadDetails"`
	Errors  []interface{} `json:"Errors"`
	Success bool          `json:"Success"`
}

type Error struct {
	Errors []struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
	} `json:"Errors"`
}
