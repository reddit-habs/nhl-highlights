// Code generated by SQLBoiler 4.14.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Highlight is an object representing the database table.
type Highlight struct {
	ID       int64      `boil:"id" json:"id" toml:"id" yaml:"id"`
	GameID   int64      `boil:"game_id" json:"game_id" toml:"game_id" yaml:"game_id"`
	MediaURL string     `boil:"media_url" json:"media_url" toml:"media_url" yaml:"media_url"`
	EventID  null.Int64 `boil:"event_id" json:"event_id,omitempty" toml:"event_id" yaml:"event_id,omitempty"`

	R *highlightR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L highlightL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var HighlightColumns = struct {
	ID       string
	GameID   string
	MediaURL string
	EventID  string
}{
	ID:       "id",
	GameID:   "game_id",
	MediaURL: "media_url",
	EventID:  "event_id",
}

var HighlightTableColumns = struct {
	ID       string
	GameID   string
	MediaURL string
	EventID  string
}{
	ID:       "highlights.id",
	GameID:   "highlights.game_id",
	MediaURL: "highlights.media_url",
	EventID:  "highlights.event_id",
}

// Generated where

type whereHelpernull_Int64 struct{ field string }

func (w whereHelpernull_Int64) EQ(x null.Int64) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Int64) NEQ(x null.Int64) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Int64) LT(x null.Int64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Int64) LTE(x null.Int64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Int64) GT(x null.Int64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Int64) GTE(x null.Int64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}
func (w whereHelpernull_Int64) IN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelpernull_Int64) NIN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

func (w whereHelpernull_Int64) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Int64) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

var HighlightWhere = struct {
	ID       whereHelperint64
	GameID   whereHelperint64
	MediaURL whereHelperstring
	EventID  whereHelpernull_Int64
}{
	ID:       whereHelperint64{field: "\"highlights\".\"id\""},
	GameID:   whereHelperint64{field: "\"highlights\".\"game_id\""},
	MediaURL: whereHelperstring{field: "\"highlights\".\"media_url\""},
	EventID:  whereHelpernull_Int64{field: "\"highlights\".\"event_id\""},
}

// HighlightRels is where relationship names are stored.
var HighlightRels = struct {
	Game string
}{
	Game: "Game",
}

// highlightR is where relationships are stored.
type highlightR struct {
	Game *Game `boil:"Game" json:"Game" toml:"Game" yaml:"Game"`
}

// NewStruct creates a new relationship struct
func (*highlightR) NewStruct() *highlightR {
	return &highlightR{}
}

func (r *highlightR) GetGame() *Game {
	if r == nil {
		return nil
	}
	return r.Game
}

// highlightL is where Load methods for each relationship are stored.
type highlightL struct{}

var (
	highlightAllColumns            = []string{"id", "game_id", "media_url", "event_id"}
	highlightColumnsWithoutDefault = []string{"game_id", "media_url"}
	highlightColumnsWithDefault    = []string{"id", "event_id"}
	highlightPrimaryKeyColumns     = []string{"id", "game_id"}
	highlightGeneratedColumns      = []string{"id"}
)

type (
	// HighlightSlice is an alias for a slice of pointers to Highlight.
	// This should almost always be used instead of []Highlight.
	HighlightSlice []*Highlight

	highlightQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	highlightType                 = reflect.TypeOf(&Highlight{})
	highlightMapping              = queries.MakeStructMapping(highlightType)
	highlightPrimaryKeyMapping, _ = queries.BindMapping(highlightType, highlightMapping, highlightPrimaryKeyColumns)
	highlightInsertCacheMut       sync.RWMutex
	highlightInsertCache          = make(map[string]insertCache)
	highlightUpdateCacheMut       sync.RWMutex
	highlightUpdateCache          = make(map[string]updateCache)
	highlightUpsertCacheMut       sync.RWMutex
	highlightUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single highlight record from the query.
func (q highlightQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Highlight, error) {
	o := &Highlight{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for highlights")
	}

	return o, nil
}

// All returns all Highlight records from the query.
func (q highlightQuery) All(ctx context.Context, exec boil.ContextExecutor) (HighlightSlice, error) {
	var o []*Highlight

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Highlight slice")
	}

	return o, nil
}

// Count returns the count of all Highlight records in the query.
func (q highlightQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count highlights rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q highlightQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if highlights exists")
	}

	return count > 0, nil
}

// Game pointed to by the foreign key.
func (o *Highlight) Game(mods ...qm.QueryMod) gameQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"game_id\" = ?", o.GameID),
	}

	queryMods = append(queryMods, mods...)

	return Games(queryMods...)
}

// LoadGame allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (highlightL) LoadGame(ctx context.Context, e boil.ContextExecutor, singular bool, maybeHighlight interface{}, mods queries.Applicator) error {
	var slice []*Highlight
	var object *Highlight

	if singular {
		var ok bool
		object, ok = maybeHighlight.(*Highlight)
		if !ok {
			object = new(Highlight)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeHighlight)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeHighlight))
			}
		}
	} else {
		s, ok := maybeHighlight.(*[]*Highlight)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeHighlight)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeHighlight))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &highlightR{}
		}
		args = append(args, object.GameID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &highlightR{}
			}

			for _, a := range args {
				if a == obj.GameID {
					continue Outer
				}
			}

			args = append(args, obj.GameID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`games`),
		qm.WhereIn(`games.game_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Game")
	}

	var resultSlice []*Game
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Game")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for games")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for games")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Game = foreign
		if foreign.R == nil {
			foreign.R = &gameR{}
		}
		foreign.R.Highlights = append(foreign.R.Highlights, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.GameID == foreign.GameID {
				local.R.Game = foreign
				if foreign.R == nil {
					foreign.R = &gameR{}
				}
				foreign.R.Highlights = append(foreign.R.Highlights, local)
				break
			}
		}
	}

	return nil
}

// SetGame of the highlight to the related item.
// Sets o.R.Game to related.
// Adds o to related.R.Highlights.
func (o *Highlight) SetGame(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Game) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"highlights\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 0, []string{"game_id"}),
		strmangle.WhereClause("\"", "\"", 0, highlightPrimaryKeyColumns),
	)
	values := []interface{}{related.GameID, o.ID, o.GameID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.GameID = related.GameID
	if o.R == nil {
		o.R = &highlightR{
			Game: related,
		}
	} else {
		o.R.Game = related
	}

	if related.R == nil {
		related.R = &gameR{
			Highlights: HighlightSlice{o},
		}
	} else {
		related.R.Highlights = append(related.R.Highlights, o)
	}

	return nil
}

// Highlights retrieves all the records using an executor.
func Highlights(mods ...qm.QueryMod) highlightQuery {
	mods = append(mods, qm.From("\"highlights\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"highlights\".*"})
	}

	return highlightQuery{q}
}

// FindHighlight retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindHighlight(ctx context.Context, exec boil.ContextExecutor, iD int64, gameID int64, selectCols ...string) (*Highlight, error) {
	highlightObj := &Highlight{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"highlights\" where \"id\"=? AND \"game_id\"=?", sel,
	)

	q := queries.Raw(query, iD, gameID)

	err := q.Bind(ctx, exec, highlightObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from highlights")
	}

	return highlightObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Highlight) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no highlights provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(highlightColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	highlightInsertCacheMut.RLock()
	cache, cached := highlightInsertCache[key]
	highlightInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			highlightAllColumns,
			highlightColumnsWithDefault,
			highlightColumnsWithoutDefault,
			nzDefaults,
		)
		wl = strmangle.SetComplement(wl, highlightGeneratedColumns)

		cache.valueMapping, err = queries.BindMapping(highlightType, highlightMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(highlightType, highlightMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"highlights\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"highlights\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into highlights")
	}

	if !cached {
		highlightInsertCacheMut.Lock()
		highlightInsertCache[key] = cache
		highlightInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Highlight.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Highlight) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	highlightUpdateCacheMut.RLock()
	cache, cached := highlightUpdateCache[key]
	highlightUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			highlightAllColumns,
			highlightPrimaryKeyColumns,
		)
		wl = strmangle.SetComplement(wl, highlightGeneratedColumns)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update highlights, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"highlights\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 0, wl),
			strmangle.WhereClause("\"", "\"", 0, highlightPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(highlightType, highlightMapping, append(wl, highlightPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update highlights row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for highlights")
	}

	if !cached {
		highlightUpdateCacheMut.Lock()
		highlightUpdateCache[key] = cache
		highlightUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q highlightQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for highlights")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for highlights")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o HighlightSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), highlightPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"highlights\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, highlightPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in highlight slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all highlight")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Highlight) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no highlights provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(highlightColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	highlightUpsertCacheMut.RLock()
	cache, cached := highlightUpsertCache[key]
	highlightUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			highlightAllColumns,
			highlightColumnsWithDefault,
			highlightColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			highlightAllColumns,
			highlightPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert highlights, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(highlightPrimaryKeyColumns))
			copy(conflict, highlightPrimaryKeyColumns)
		}
		cache.query = buildUpsertQuerySQLite(dialect, "\"highlights\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(highlightType, highlightMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(highlightType, highlightMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert highlights")
	}

	if !cached {
		highlightUpsertCacheMut.Lock()
		highlightUpsertCache[key] = cache
		highlightUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Highlight record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Highlight) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Highlight provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), highlightPrimaryKeyMapping)
	sql := "DELETE FROM \"highlights\" WHERE \"id\"=? AND \"game_id\"=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from highlights")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for highlights")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q highlightQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no highlightQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from highlights")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for highlights")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o HighlightSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), highlightPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"highlights\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, highlightPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from highlight slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for highlights")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Highlight) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindHighlight(ctx, exec, o.ID, o.GameID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *HighlightSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := HighlightSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), highlightPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"highlights\".* FROM \"highlights\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, highlightPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in HighlightSlice")
	}

	*o = slice

	return nil
}

// HighlightExists checks if the Highlight row exists.
func HighlightExists(ctx context.Context, exec boil.ContextExecutor, iD int64, gameID int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"highlights\" where \"id\"=? AND \"game_id\"=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD, gameID)
	}
	row := exec.QueryRowContext(ctx, sql, iD, gameID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if highlights exists")
	}

	return exists, nil
}

// Exists checks if the Highlight row exists.
func (o *Highlight) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return HighlightExists(ctx, exec, o.ID, o.GameID)
}
