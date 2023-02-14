package dto

//type QuoteDetailsDto struct {
//	BanyanRequest struct {
//		CustmId string `json:"custmId"`
//		Items   []struct {
//			Class    string `json:"class"`
//			Weight   string `json:"weight"`
//			Length   string `json:"length"`
//			Width    string `json:"width"`
//			Height   string `json:"height"`
//			Type     string `json:"type"`
//			Quantity string `json:"quantity"`
//		} `json:"items"`
//		Events []struct {
//			Date    string `json:"date"`
//			City    string `json:"city"`
//			State   string `json:"state"`
//			Zip     string `json:"zip"`
//			Country string `json:"country"`
//		} `json:"events"`
//	} `json:"BanyanRequest"`
//	BanyanResponse struct {
//		MercuryResponseDto struct {
//			XMLName struct {
//				Space string `json:"Space"`
//				Local string `json:"Local"`
//			} `json:"XMLName"`
//			SpecVersion   string `json:"SpecVersion"`
//			StatusCode    string `json:"StatusCode"`
//			StatusMessage string `json:"StatusMessage"`
//			PriceSheets   struct {
//				PriceSheet []struct {
//					Type             string `json:"Type"`
//					ChargeModel      string `json:"ChargeModel"`
//					IsSelected       string `json:"IsSelected"`
//					IsAllocated      string `json:"IsAllocated"`
//					CurrencyCode     string `json:"CurrencyCode"`
//					CreateDate       string `json:"CreateDate"`
//					InternalId       string `json:"InternalId"`
//					AccessorialTotal string `json:"AccessorialTotal"`
//					SubTotal         string `json:"SubTotal"`
//					Total            string `json:"Total"`
//					ContractId       string `json:"ContractId"`
//					ContractName     string `json:"ContractName"`
//					CarrierId        string `json:"CarrierId"`
//					CarrierName      string `json:"CarrierName"`
//					SCAC             string `json:"SCAC"`
//					Mode             string `json:"Mode"`
//					Service          string `json:"Service"`
//					ServiceDays      string `json:"ServiceDays"`
//					Distance         string `json:"Distance"`
//					ID               string `json:"ID"`
//					InsuranceTypes   struct {
//						Insurance []struct {
//							Type           string `json:"Type"`
//							Amount         string `json:"Amount"`
//							Company        string `json:"Company"`
//							ExpirationDate string `json:"ExpirationDate"`
//							ContactName    string `json:"ContactName"`
//							ContactPhone   string `json:"ContactPhone"`
//						} `json:"Insurance"`
//					} `json:"InsuranceTypes"`
//					Address struct {
//						Type          string `json:"Type"`
//						IsResidential string `json:"IsResidential"`
//						IsPrimary     string `json:"IsPrimary"`
//						LocationCode  string `json:"LocationCode"`
//						Alias         string `json:"Alias"`
//						Name          string `json:"Name"`
//						AddrLine1     string `json:"AddrLine1"`
//						AddrLine2     string `json:"AddrLine2"`
//						City          string `json:"City"`
//						StateProvince string `json:"StateProvince"`
//						PostalCode    string `json:"PostalCode"`
//						CountryCode   string `json:"CountryCode"`
//						GeoLoc        struct {
//							LatDegrees    string `json:"LatDegrees"`
//							LatDirection  string `json:"LatDirection"`
//							LongDegrees   string `json:"LongDegrees"`
//							LongDirection string `json:"LongDirection"`
//						} `json:"GeoLoc"`
//						Contacts struct {
//							Contact struct {
//								Type           string `json:"Type"`
//								Oid            string `json:"Oid"`
//								Name           string `json:"Name"`
//								ContactMethods struct {
//									ContactMethod []struct {
//										SequenceNum string `json:"SequenceNum"`
//										Type        string `json:"Type"`
//									} `json:"ContactMethod"`
//								} `json:"ContactMethods"`
//							} `json:"Contact"`
//						} `json:"Contacts"`
//						Comments string `json:"Comments"`
//					} `json:"Address"`
//					ExpectedDeliveryDate string `json:"ExpectedDeliveryDate"`
//					ReasonCode           string `json:"ReasonCode"`
//					Status               string `json:"Status"`
//					LaneID               string `json:"LaneID"`
//					Zone                 string `json:"Zone"`
//					RouteGuidePriority   string `json:"RouteGuidePriority"`
//					CarrierLocationOid   string `json:"CarrierLocationOid"`
//					OriginService        string `json:"OriginService"`
//					DestinationService   string `json:"DestinationService"`
//					Charges              struct {
//						Charge []struct {
//							SequenceNum     string `json:"SequenceNum"`
//							Type            string `json:"Type"`
//							ItemGroupId     string `json:"ItemGroupId"`
//							Description     string `json:"Description"`
//							EdiCode         string `json:"EdiCode"`
//							Amount          string `json:"Amount"`
//							Rate            string `json:"Rate"`
//							RateQualifier   string `json:"RateQualifier"`
//							Quantity        string `json:"Quantity"`
//							Weight          string `json:"Weight"`
//							DimWeight       string `json:"DimWeight"`
//							FreightClass    string `json:"FreightClass"`
//							FakFreightClass string `json:"FakFreightClass"`
//							IsMin           string `json:"IsMin"`
//							IsMax           string `json:"IsMax"`
//							IsNontaxable    string `json:"IsNontaxable"`
//						} `json:"Charge"`
//					} `json:"Charges"`
//					Comments         string `json:"Comments"`
//					QuoteInformation struct {
//						QuoteNumber string `json:"QuoteNumber"`
//						Date        struct {
//							Type string `json:"Type"`
//						} `json:"Date"`
//						QuoteBy    string `json:"QuoteBy"`
//						QuotePhone string `json:"QuotePhone"`
//						QuoteFax   string `json:"QuoteFax"`
//						QuoteEmail string `json:"QuoteEmail"`
//					} `json:"QuoteInformation"`
//					AssociatedCarrierPricesheet struct {
//						PriceSheet struct {
//							Type             string `json:"Type"`
//							ChargeModel      string `json:"ChargeModel"`
//							IsSelected       string `json:"IsSelected"`
//							IsAllocated      string `json:"IsAllocated"`
//							CurrencyCode     string `json:"CurrencyCode"`
//							CreateDate       string `json:"CreateDate"`
//							InternalId       string `json:"InternalId"`
//							AccessorialTotal string `json:"AccessorialTotal"`
//							SubTotal         string `json:"SubTotal"`
//							Total            string `json:"Total"`
//							ContractId       string `json:"ContractId"`
//							ContractName     string `json:"ContractName"`
//							CarrierId        string `json:"CarrierId"`
//							CarrierName      string `json:"CarrierName"`
//							SCAC             string `json:"SCAC"`
//							Mode             string `json:"Mode"`
//							Service          string `json:"Service"`
//							ServiceDays      string `json:"ServiceDays"`
//							Distance         string `json:"Distance"`
//							ID               string `json:"ID"`
//							InsuranceTypes   struct {
//								Insurance []struct {
//									Type           string `json:"Type"`
//									Amount         string `json:"Amount"`
//									Company        string `json:"Company"`
//									ExpirationDate string `json:"ExpirationDate"`
//									ContactName    string `json:"ContactName"`
//									ContactPhone   string `json:"ContactPhone"`
//								} `json:"Insurance"`
//							} `json:"InsuranceTypes"`
//							Address struct {
//								Type          string `json:"Type"`
//								IsResidential string `json:"IsResidential"`
//								IsPrimary     string `json:"IsPrimary"`
//								LocationCode  string `json:"LocationCode"`
//								Alias         string `json:"Alias"`
//								Name          string `json:"Name"`
//								AddrLine1     string `json:"AddrLine1"`
//								AddrLine2     string `json:"AddrLine2"`
//								City          string `json:"City"`
//								StateProvince string `json:"StateProvince"`
//								PostalCode    string `json:"PostalCode"`
//								CountryCode   string `json:"CountryCode"`
//								GeoLoc        struct {
//									LatDegrees    string `json:"LatDegrees"`
//									LatDirection  string `json:"LatDirection"`
//									LongDegrees   string `json:"LongDegrees"`
//									LongDirection string `json:"LongDirection"`
//								} `json:"GeoLoc"`
//								Contacts struct {
//									Contact struct {
//										Type           string `json:"Type"`
//										Oid            string `json:"Oid"`
//										Name           string `json:"Name"`
//										ContactMethods struct {
//											ContactMethod []struct {
//												SequenceNum string `json:"SequenceNum"`
//												Type        string `json:"Type"`
//											} `json:"ContactMethod"`
//										} `json:"ContactMethods"`
//									} `json:"Contact"`
//								} `json:"Contacts"`
//								Comments string `json:"Comments"`
//							} `json:"Address"`
//							ExpectedDeliveryDate string `json:"ExpectedDeliveryDate"`
//							ReasonCode           string `json:"ReasonCode"`
//							Status               string `json:"Status"`
//							LaneID               string `json:"LaneID"`
//							Zone                 string `json:"Zone"`
//							RouteGuidePriority   string `json:"RouteGuidePriority"`
//							CarrierLocationOid   string `json:"CarrierLocationOid"`
//							OriginService        string `json:"OriginService"`
//							DestinationService   string `json:"DestinationService"`
//							Charges              struct {
//								Charge []struct {
//									SequenceNum     string `json:"SequenceNum"`
//									Type            string `json:"Type"`
//									ItemGroupId     string `json:"ItemGroupId"`
//									Description     string `json:"Description"`
//									EdiCode         string `json:"EdiCode"`
//									Amount          string `json:"Amount"`
//									Rate            string `json:"Rate"`
//									RateQualifier   string `json:"RateQualifier"`
//									Quantity        string `json:"Quantity"`
//									Weight          string `json:"Weight"`
//									DimWeight       string `json:"DimWeight"`
//									FreightClass    string `json:"FreightClass"`
//									FakFreightClass string `json:"FakFreightClass"`
//									IsMin           string `json:"IsMin"`
//									IsMax           string `json:"IsMax"`
//									IsNontaxable    string `json:"IsNontaxable"`
//								} `json:"Charge"`
//							} `json:"Charges"`
//							Comments         string `json:"Comments"`
//							QuoteInformation struct {
//								QuoteNumber string `json:"QuoteNumber"`
//								Date        struct {
//									Type string `json:"Type"`
//								} `json:"Date"`
//								QuoteBy    string `json:"QuoteBy"`
//								QuotePhone string `json:"QuotePhone"`
//								QuoteFax   string `json:"QuoteFax"`
//								QuoteEmail string `json:"QuoteEmail"`
//							} `json:"QuoteInformation"`
//						} `json:"PriceSheet"`
//					} `json:"AssociatedCarrierPricesheet"`
//				} `json:"PriceSheet"`
//			} `json:"PriceSheets"`
//		} `json:"MercuryResponseDto"`
//	} `json:"BanyanResponse"`
//}
type QuoteDetailsDto struct {
	Request struct {
		CustmId string `json:"custmId"`
		Items   []struct {
			Class    string `json:"class"`
			Weight   string `json:"weight"`
			Length   string `json:"length"`
			Width    string `json:"width"`
			Height   string `json:"height"`
			Type     string `json:"type"`
			Quantity string `json:"quantity"`
		} `json:"items"`
		Events []struct {
			Date    string `json:"date"`
			City    string `json:"city"`
			State   string `json:"state"`
			Zip     string `json:"zip"`
			Country string `json:"country"`
		} `json:"events"`
	} `json:"Request"`
	Response struct {
		MercuryResponseDto struct {
			XMLName struct {
				Space string `json:"Space"`
				Local string `json:"Local"`
			} `json:"XMLName"`
			SpecVersion   string `json:"SpecVersion"`
			StatusCode    string `json:"StatusCode"`
			StatusMessage string `json:"StatusMessage"`
			PriceSheets   struct {
				PriceSheet []struct {
					Type             string `json:"Type"`
					ChargeModel      string `json:"ChargeModel"`
					IsSelected       string `json:"IsSelected"`
					IsAllocated      string `json:"IsAllocated"`
					CurrencyCode     string `json:"CurrencyCode"`
					CreateDate       string `json:"CreateDate"`
					InternalId       string `json:"InternalId"`
					AccessorialTotal string `json:"AccessorialTotal"`
					SubTotal         string `json:"SubTotal"`
					Total            string `json:"Total"`
					ContractId       string `json:"ContractId"`
					ContractName     string `json:"ContractName"`
					CarrierId        string `json:"CarrierId"`
					CarrierName      string `json:"CarrierName"`
					SCAC             string `json:"SCAC"`
					Mode             string `json:"Mode"`
					Service          string `json:"Service"`
					ServiceDays      string `json:"ServiceDays"`
					Distance         string `json:"Distance"`
					ID               string `json:"ID"`
					InsuranceTypes   struct {
						Insurance []struct {
							Type           string `json:"Type"`
							Amount         string `json:"Amount"`
							Company        string `json:"Company"`
							ExpirationDate string `json:"ExpirationDate"`
							ContactName    string `json:"ContactName"`
							ContactPhone   string `json:"ContactPhone"`
						} `json:"Insurance"`
					} `json:"InsuranceTypes"`
					Address struct {
						Type          string `json:"Type"`
						IsResidential string `json:"IsResidential"`
						IsPrimary     string `json:"IsPrimary"`
						LocationCode  string `json:"LocationCode"`
						Alias         string `json:"Alias"`
						Name          string `json:"Name"`
						AddrLine1     string `json:"AddrLine1"`
						AddrLine2     string `json:"AddrLine2"`
						City          string `json:"City"`
						StateProvince string `json:"StateProvince"`
						PostalCode    string `json:"PostalCode"`
						CountryCode   string `json:"CountryCode"`
						GeoLoc        struct {
							LatDegrees    string `json:"LatDegrees"`
							LatDirection  string `json:"LatDirection"`
							LongDegrees   string `json:"LongDegrees"`
							LongDirection string `json:"LongDirection"`
						} `json:"GeoLoc"`
						Contacts struct {
							Contact struct {
								Type           string `json:"Type"`
								Oid            string `json:"Oid"`
								Name           string `json:"Name"`
								ContactMethods struct {
									ContactMethod []struct {
										SequenceNum string `json:"SequenceNum"`
										Type        string `json:"Type"`
									} `json:"ContactMethod"`
								} `json:"ContactMethods"`
							} `json:"Contact"`
						} `json:"Contacts"`
						Comments string `json:"Comments"`
					} `json:"Address"`
					ExpectedDeliveryDate string `json:"ExpectedDeliveryDate"`
					ReasonCode           string `json:"ReasonCode"`
					Status               string `json:"Status"`
					LaneID               string `json:"LaneID"`
					Zone                 string `json:"Zone"`
					RouteGuidePriority   string `json:"RouteGuidePriority"`
					CarrierLocationOid   string `json:"CarrierLocationOid"`
					OriginService        string `json:"OriginService"`
					DestinationService   string `json:"DestinationService"`
					Charges              struct {
						Charge []struct {
							SequenceNum     string `json:"SequenceNum"`
							Type            string `json:"Type"`
							ItemGroupId     string `json:"ItemGroupId"`
							Description     string `json:"Description"`
							EdiCode         string `json:"EdiCode"`
							Amount          string `json:"Amount"`
							Rate            string `json:"Rate"`
							RateQualifier   string `json:"RateQualifier"`
							Quantity        string `json:"Quantity"`
							Weight          string `json:"Weight"`
							DimWeight       string `json:"DimWeight"`
							FreightClass    string `json:"FreightClass"`
							FakFreightClass string `json:"FakFreightClass"`
							IsMin           string `json:"IsMin"`
							IsMax           string `json:"IsMax"`
							IsNontaxable    string `json:"IsNontaxable"`
						} `json:"Charge"`
					} `json:"Charges"`
					Comments         string `json:"Comments"`
					QuoteInformation struct {
						QuoteNumber string `json:"QuoteNumber"`
						Date        struct {
							Type string `json:"Type"`
						} `json:"Date"`
						QuoteBy    string `json:"QuoteBy"`
						QuotePhone string `json:"QuotePhone"`
						QuoteFax   string `json:"QuoteFax"`
						QuoteEmail string `json:"QuoteEmail"`
					} `json:"QuoteInformation"`
					AssociatedCarrierPricesheet struct {
						PriceSheet struct {
							Type             string `json:"Type"`
							ChargeModel      string `json:"ChargeModel"`
							IsSelected       string `json:"IsSelected"`
							IsAllocated      string `json:"IsAllocated"`
							CurrencyCode     string `json:"CurrencyCode"`
							CreateDate       string `json:"CreateDate"`
							InternalId       string `json:"InternalId"`
							AccessorialTotal string `json:"AccessorialTotal"`
							SubTotal         string `json:"SubTotal"`
							Total            string `json:"Total"`
							ContractId       string `json:"ContractId"`
							ContractName     string `json:"ContractName"`
							CarrierId        string `json:"CarrierId"`
							CarrierName      string `json:"CarrierName"`
							SCAC             string `json:"SCAC"`
							Mode             string `json:"Mode"`
							Service          string `json:"Service"`
							ServiceDays      string `json:"ServiceDays"`
							Distance         string `json:"Distance"`
							ID               string `json:"ID"`
							InsuranceTypes   struct {
								Insurance []struct {
									Type           string `json:"Type"`
									Amount         string `json:"Amount"`
									Company        string `json:"Company"`
									ExpirationDate string `json:"ExpirationDate"`
									ContactName    string `json:"ContactName"`
									ContactPhone   string `json:"ContactPhone"`
								} `json:"Insurance"`
							} `json:"InsuranceTypes"`
							Address struct {
								Type          string `json:"Type"`
								IsResidential string `json:"IsResidential"`
								IsPrimary     string `json:"IsPrimary"`
								LocationCode  string `json:"LocationCode"`
								Alias         string `json:"Alias"`
								Name          string `json:"Name"`
								AddrLine1     string `json:"AddrLine1"`
								AddrLine2     string `json:"AddrLine2"`
								City          string `json:"City"`
								StateProvince string `json:"StateProvince"`
								PostalCode    string `json:"PostalCode"`
								CountryCode   string `json:"CountryCode"`
								GeoLoc        struct {
									LatDegrees    string `json:"LatDegrees"`
									LatDirection  string `json:"LatDirection"`
									LongDegrees   string `json:"LongDegrees"`
									LongDirection string `json:"LongDirection"`
								} `json:"GeoLoc"`
								Contacts struct {
									Contact struct {
										Type           string `json:"Type"`
										Oid            string `json:"Oid"`
										Name           string `json:"Name"`
										ContactMethods struct {
											ContactMethod []struct {
												SequenceNum string `json:"SequenceNum"`
												Type        string `json:"Type"`
											} `json:"ContactMethod"`
										} `json:"ContactMethods"`
									} `json:"Contact"`
								} `json:"Contacts"`
								Comments string `json:"Comments"`
							} `json:"Address"`
							ExpectedDeliveryDate string `json:"ExpectedDeliveryDate"`
							ReasonCode           string `json:"ReasonCode"`
							Status               string `json:"Status"`
							LaneID               string `json:"LaneID"`
							Zone                 string `json:"Zone"`
							RouteGuidePriority   string `json:"RouteGuidePriority"`
							CarrierLocationOid   string `json:"CarrierLocationOid"`
							OriginService        string `json:"OriginService"`
							DestinationService   string `json:"DestinationService"`
							Charges              struct {
								Charge []struct {
									SequenceNum     string `json:"SequenceNum"`
									Type            string `json:"Type"`
									ItemGroupId     string `json:"ItemGroupId"`
									Description     string `json:"Description"`
									EdiCode         string `json:"EdiCode"`
									Amount          string `json:"Amount"`
									Rate            string `json:"Rate"`
									RateQualifier   string `json:"RateQualifier"`
									Quantity        string `json:"Quantity"`
									Weight          string `json:"Weight"`
									DimWeight       string `json:"DimWeight"`
									FreightClass    string `json:"FreightClass"`
									FakFreightClass string `json:"FakFreightClass"`
									IsMin           string `json:"IsMin"`
									IsMax           string `json:"IsMax"`
									IsNontaxable    string `json:"IsNontaxable"`
								} `json:"Charge"`
							} `json:"Charges"`
							Comments         string `json:"Comments"`
							QuoteInformation struct {
								QuoteNumber string `json:"QuoteNumber"`
								Date        struct {
									Type string `json:"Type"`
								} `json:"Date"`
								QuoteBy    string `json:"QuoteBy"`
								QuotePhone string `json:"QuotePhone"`
								QuoteFax   string `json:"QuoteFax"`
								QuoteEmail string `json:"QuoteEmail"`
							} `json:"QuoteInformation"`
						} `json:"PriceSheet"`
					} `json:"AssociatedCarrierPricesheet"`
				} `json:"PriceSheet"`
			} `json:"PriceSheets"`
		} `json:"MercuryResponseDto"`
	} `json:"Response"`
}
