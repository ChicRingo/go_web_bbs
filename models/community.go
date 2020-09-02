package models

import "time"

type Community struct {
	ID   int64  `db:"community_id" json:"id"`
	Name string `db:"community_name" json:"name"`
}

type CommunityDetail struct {
	ID           int64     `db:"community_id" json:"id"`
	Name         string    `db:"community_name" json:"name"`
	Introduction string    `db:"introduction" json:"introduction,omitempty"`
	CreateTime   time.Time `db:"create_time" json:"create_time"`
	UpdateTime   time.Time `db:"update_time" json:"update_time"`
}
