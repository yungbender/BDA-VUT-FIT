package crawler_collector

import (
	"bda/models"
	"bda/types"
	"database/sql"
	"encoding/binary"
	"net"
	"time"

	"gorm.io/gorm"
)

func SaveAddr(db *gorm.DB, addr types.AddrChanMsg) {
	newIp := net.IP(addr.Addr.IpAddr[:])
	var node models.Node
	dirty := false
	db.FirstOrInit(&node, models.Node{Ip: newIp.String(),
		Port: binary.BigEndian.Uint16(addr.Addr.Port[:])})

	if !node.FromIp.Valid {
		node.FromIp = sql.NullString{
			String: addr.NodeIp.String(),
			Valid:  true,
		}
		dirty = true
	}

	if !node.Discovered.Valid {
		node.Discovered = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		dirty = true
	}

	if !node.Connected.Valid || (node.Connected.Valid && node.Connected.Time.After(time.Unix(int64(addr.Addr.Time), 0))) {
		node.Connected = sql.NullTime{
			Time:  time.Unix(int64(addr.Addr.Time), 0),
			Valid: true,
		}
		dirty = true
		// if we got ADDR about node which was connected EARLIER than we have in the databse, override it with older date
	}

	if dirty {
		db.Save(&node)
	}
}

func SaveVersion(db *gorm.DB, version types.VersionChanMsg) {
	var node models.Node
	dirty := false
	db.FirstOrInit(&node, models.Node{Ip: version.NodeIp.String(), Port: version.NodePort})

	// Always update useragent to have freshest info
	if !node.UserAgent.Valid || (node.UserAgent.Valid && node.UserAgent.String != version.UserAgent) {
		node.UserAgent = sql.NullString{
			String: version.UserAgent,
			Valid:  true,
		}
		dirty = true
	}

	// if time is somehow further than current time, update it to current time, just in case
	if !node.Connected.Valid || (node.Connected.Valid && node.Connected.Time.After(time.Now())) {
		node.Connected = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		dirty = true
	}

	if dirty {
		db.Save(&node)
	}
}
