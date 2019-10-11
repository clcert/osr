package models


// DefaultModels defines the default model list.
var DefaultModels = ModelsList{
	Name:   "default",
	Models: []Model{
		// Core Models
		TaskModel,
		// Basic Metainfo
		SourceModel,
		ContactModel,
		CountryModel,
		ASNModel,
		DomainModel,
		PortModel,
		// Advanced Metainfo
		BlacklistedSubnetModel,
		SubnetCountryModel,
		SubnetASNModel,
		SubnetContactModel,
		DomainRankingModel,
		DomainCategoryModel,
		DomainToCategoryModel,
		ReportTypeModel,
		// Scans
		CertificateModel,
		DarknetPacketModel,
		DnsRRModel,
		IpAsnCountryModel, // Intermediate Tabld
		PortScanModel,
		ReportEntryModel,
	},
}