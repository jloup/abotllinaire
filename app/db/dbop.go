package db

import (
	"github.com/jloup/utils"
	"gopkg.in/mgo.v2"
)

type DbOp struct {
	c     *mgo.Collection
	query *mgo.Query

	err error
}

func (d *DbOp) Err() error {
	if d.err == nil {
		return nil
	}

	if d.err == mgo.ErrNotFound {
		return utils.Error{NotFound, d.err.Error()}
	}

	if mgo.IsDup(d.err) {
		return utils.Error{IsDup, d.err.Error()}
	}

	return utils.Error{DbError, d.err.Error()}
}

func (d *DbOp) Find(query interface{}) *DbOp {
	if d.err != nil {
		return d
	}

	d.query = d.c.Find(query)
	return d
}

func (d *DbOp) Count() (int, error) {
	var n int
	if d.err != nil {
		return 0, nil
	}

	n, d.err = d.query.Count()

	return n, d.err
}

func (d *DbOp) Sort(fields ...string) *DbOp {
	if d.err != nil {
		return d
	}

	d.query = d.query.Sort(fields...)
	return d
}

func (d *DbOp) Limit(n int) *DbOp {
	if d.err != nil {
		return d
	}

	d.query = d.query.Limit(n)
	return d
}

func (d *DbOp) All(result interface{}) {
	if d.err != nil {
		return
	}

	d.err = d.query.All(result)
}

func (d *DbOp) Iter() *mgo.Iter {
	if d.err != nil {
		return nil
	}

	return d.query.Iter()
}

func (d *DbOp) One(result interface{}) {
	if d.err != nil {
		return
	}

	d.err = d.query.One(result)
}

func (d *DbOp) Insert(docs ...interface{}) {
	if d.err != nil {
		return
	}

	d.err = d.c.Insert(docs...)
}

func (d *DbOp) Update(selector, update interface{}) {
	if d.err != nil {
		return
	}

	d.err = d.c.Update(selector, update)
}
func (d *DbOp) Remove(selector interface{}) {
	if d.err != nil {
		return
	}

	d.err = d.c.Remove(selector)
}

func (d *DbOp) RemoveId(id interface{}) {
	if d.err != nil {
		return
	}

	d.err = d.c.RemoveId(id)
}

func (d *DbOp) UpdateAll(selector, update interface{}) *mgo.ChangeInfo {
	if d.err != nil {
		return nil
	}

	var info *mgo.ChangeInfo
	info, d.err = d.c.UpdateAll(selector, update)

	return info
}

func (d *DbOp) RemoveAll(selector interface{}) *mgo.ChangeInfo {
	if d.err != nil {
		return nil
	}

	var info *mgo.ChangeInfo
	info, d.err = d.c.RemoveAll(selector)

	return info
}
