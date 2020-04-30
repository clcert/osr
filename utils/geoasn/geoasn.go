package geoasn

import (
	"fmt"
	"github.com/clcert/osr/models"
	"github.com/go-pg/pg"
	"net"
)

// Classifier allows to clasify IPs based on Geographical and ASN metadata
type Classifier struct {
	GeoTree *IPNetTree
	ASNTree *IPNetTree
}

// NewClassifier returns a new classifier using data from the latest import from a specific source
func NewClassifier(db *pg.DB, source models.DataSourceID) (classifier *Classifier, err error) {
	geoTree, err := getGeoTree(db, source)
	if err != nil {
		return
	}
	asnTree, err := getASNTree(db, source)
	if err != nil {
		return
	}
	classifier = &Classifier{
		GeoTree: geoTree,
		ASNTree: asnTree,
	}
	return
}

// GetGeoASN returns the Geogrhaphical and ASN data for an specific IP
func (classifier *Classifier) GetGeoASN(ip net.IP) (int, int, error) {
	geoData, ok := classifier.GeoTree.GetIPData(ip)
	if !ok {
		return 0, 0, fmt.Errorf("cannot get IP Geodata for IP %s", ip)
	}
	asnData, ok := classifier.ASNTree.GetIPData(ip)
	if !ok {
		return 0, 0, fmt.Errorf("cannot get IP ASN data for IP %s", ip)
	}
	return geoData.(int), asnData.(int), nil
}

func getGeoTree(db *pg.DB, source models.DataSourceID) (tree *IPNetTree, err error) {
	taskID, err := models.LatestModelTaskID(db, &models.SubnetCountry{})
	if err != nil {
		return
	}
	subnetList := make([]*models.SubnetCountry, 0)
	if err = db.Model(&models.SubnetCountry{}).
		Column("subnet", "country_geoname_id").
		Where("task_id = ? and source_id = ?", taskID, source).
		Order("subnet ASC").Select(&subnetList); err != nil {
		return
	}

	tree = NewIPNetTree()
	for _, elem := range subnetList {
		tree.AddNode(&IPNetNode{
			IPNet: elem.Subnet,
			Value: elem.CountryGeonameId,
		})
	}
	return
}

func getASNTree(db *pg.DB, source models.DataSourceID) (tree *IPNetTree, err error) {
	taskID, err := models.LatestModelTaskID(db, &models.SubnetASN{})
	if err != nil {
		return
	}
	subnetList := make([]*models.SubnetASN, 0)
	if err = db.Model(&models.SubnetASN{}).
		Column("subnet", "asn_id").
		Where("task_id = ? and source_id = ?", taskID, source).
		Order("subnet ASC").Select(&subnetList); err != nil {
		return
	}

	tree = NewIPNetTree()
	for _, elem := range subnetList {
		tree.AddNode(&IPNetNode{
			IPNet: elem.Subnet,
			Value: elem.AsnID,
		})
	}
	return
}
