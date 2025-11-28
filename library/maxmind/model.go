package maxmind

import "net/netip"

// The Result struct corresponds to the data in the GeoIP2/GeoLite2 City databases.
type Result struct {
	Traits             CityTraits         `json:"traits,omitzero"              maxminddb:"traits"`              // Traits contains various traits associated with the IP address
	Postal             CityPostal         `json:"postal,omitzero"              maxminddb:"postal"`              // Postal contains data for the postal record associated with the IP address
	Continent          Continent          `json:"continent,omitzero"           maxminddb:"continent"`           // Continent contains data for the continent record associated with the IP address
	City               CityRecord         `json:"city,omitzero"                maxminddb:"city"`                // City contains data for the city record associated with the IP address
	Subdivisions       []CitySubdivision  `json:"subdivisions,omitzero"        maxminddb:"subdivisions"`        // Subdivisions contains data for the subdivisions associated with the IP address. The subdivisions array is ordered from largest to smallest. For instance, the response for Oxford in the United Kingdom would have England as the first element and Oxfordshire as the second element.
	RepresentedCountry RepresentedCountry `json:"represented_country,omitzero" maxminddb:"represented_country"` // RepresentedCountry contains data for the represented country associated with the IP address. The represented country is the country represented by something like a military base or embassy.
	Country            CountryRecord      `json:"country,omitzero"             maxminddb:"country"`             // Country contains data for the country record associated with the IP address. This record represents the country where MaxMind believes the IP is located.
	RegisteredCountry  CountryRecord      `json:"registered_country,omitzero"  maxminddb:"registered_country"`  // RegisteredCountry contains data for the registered country associated with the IP address. This record represents the country where the ISP has registered the IP block and may differ from the user's country.
	Location           Location           `json:"location,omitzero"            maxminddb:"location"`            // Location contains data for the location record associated with the IP address
}

// Names contains localized names for geographic entities.
type Names struct {
	German              string `json:"de,omitzero"    maxminddb:"de"`    // German localized name
	English             string `json:"en,omitzero"    maxminddb:"en"`    // English localized name
	Spanish             string `json:"es,omitzero"    maxminddb:"es"`    // Spanish localized name
	French              string `json:"fr,omitzero"    maxminddb:"fr"`    // French localized name
	Japanese            string `json:"ja,omitzero"    maxminddb:"ja"`    // Japanese localized name
	BrazilianPortuguese string `json:"pt-BR,omitzero" maxminddb:"pt-BR"` // BrazilianPortuguese localized name (pt-BR)
	Russian             string `json:"ru,omitzero"    maxminddb:"ru"`    // Russian localized name
	SimplifiedChinese   string `json:"zh-CN,omitzero" maxminddb:"zh-CN"` // SimplifiedChinese localized name (zh-CN)
}

// Continent contains data for the continent record associated with an IP address.
type Continent struct {
	Names     Names  `json:"names,omitzero"      maxminddb:"names"`      // Names contains localized names for the continent
	Code      string `json:"code,omitzero"       maxminddb:"code"`       // Code is a two character continent code like "NA" (North America) or "OC" (Oceania)
	GeoNameID uint   `json:"geoname_id,omitzero" maxminddb:"geoname_id"` // GeoNameID for the continent
}

// Location contains data for the location record associated with an IP address.
type Location struct {
	// Latitude is the approximate latitude of the location associated with
	// the IP address. This value is not precise and should not be used to
	// identify a particular address or household. Will be nil if not present
	// in the database.
	Latitude *float64 `json:"latitude,omitzero"        maxminddb:"latitude"`
	// Longitude is the approximate longitude of the location associated with
	// the IP address. This value is not precise and should not be used to
	// identify a particular address or household. Will be nil if not present
	// in the database.
	Longitude *float64 `json:"longitude,omitzero"       maxminddb:"longitude"`
	// TimeZone is the time zone associated with location, as specified by
	// the IANA Time Zone Database (e.g., "America/New_York")
	TimeZone string `json:"time_zone,omitzero"       maxminddb:"time_zone"`
	// MetroCode is a metro code for targeting advertisements.
	//
	// Deprecated: Metro codes are no longer maintained and should not be used.
	MetroCode uint `json:"metro_code,omitzero"      maxminddb:"metro_code"`
	// AccuracyRadius is the approximate accuracy radius in kilometers around
	// the latitude and longitude. This is the radius where we have a 67%
	// confidence that the device using the IP address resides within the
	// circle.
	AccuracyRadius uint16 `json:"accuracy_radius,omitzero" maxminddb:"accuracy_radius"`
}

// RepresentedCountry contains data for the represented country associated
// with an IP address. The represented country is the country represented
// by something like a military base or embassy.
type RepresentedCountry struct {
	Names             Names  `json:"names,omitzero"                maxminddb:"names"`                // Names contains localized names for the represented country
	ISOCode           string `json:"iso_code,omitzero"             maxminddb:"iso_code"`             // ISOCode is the two-character ISO 3166-1 alpha code for the represented country. See https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
	Type              string `json:"type,omitzero"                 maxminddb:"type"`                 // Type is a string indicating the type of entity that is representing the country. Currently this is only "military" but may expand in the future.
	GeoNameID         uint   `json:"geoname_id,omitzero"           maxminddb:"geoname_id"`           // GeoNameID for the represented country
	IsInEuropeanUnion bool   `json:"is_in_european_union,omitzero" maxminddb:"is_in_european_union"` // IsInEuropeanUnion is true if the represented country is a member state of the European Union
}

// CityRecord contains city data for City database records.
type CityRecord struct {
	Names     Names `json:"names,omitzero"      maxminddb:"names"`      // Names contains localized names for the city
	GeoNameID uint  `json:"geoname_id,omitzero" maxminddb:"geoname_id"` // GeoNameID for the city
}

// CityPostal contains postal data for City database records.
type CityPostal struct {
	Code string `json:"code,omitzero" maxminddb:"code"` // Code is the postal code of the location. Postal codes are not available for all countries. In some countries, this will only contain part of the postal code.
}

// CitySubdivision contains subdivision data for City database records.
type CitySubdivision struct {
	Names     Names  `json:"names,omitzero"      maxminddb:"names"`
	ISOCode   string `json:"iso_code,omitzero"   maxminddb:"iso_code"`
	GeoNameID uint   `json:"geoname_id,omitzero" maxminddb:"geoname_id"`
}

// CountryRecord contains country data for City and Country database records.
type CountryRecord struct {
	Names             Names  `json:"names,omitzero"                maxminddb:"names"`                // Names contains localized names for the country
	ISOCode           string `json:"iso_code,omitzero"             maxminddb:"iso_code"`             // ISOCode is the two-character ISO 3166-1 alpha code for the country. See https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
	GeoNameID         uint   `json:"geoname_id,omitzero"           maxminddb:"geoname_id"`           // GeoNameID for the country
	IsInEuropeanUnion bool   `json:"is_in_european_union,omitzero" maxminddb:"is_in_european_union"` // IsInEuropeanUnion is true if the country is a member state of the European Union
}

// CityTraits contains traits data for City database records.
type CityTraits struct {
	IPAddress netip.Addr   `json:"ip_address,omitzero"`                        // IPAddress is the IP address used during the lookup
	Network   netip.Prefix `json:"network,omitzero"`                           // Network is the network prefix for this record. This is the largest network where all of the fields besides IPAddress have the same value.
	IsAnycast bool         `json:"is_anycast,omitzero" maxminddb:"is_anycast"` // IsAnycast is true if the IP address belongs to an anycast network. See https://en.wikipedia.org/wiki/Anycast
}
