package donelog

// DoneLog is the aggregate root that represents a single record of work done.
type DoneLog struct {
	id         DoneLogID
	title      Title
	trackID    TrackID
	categoryID CategoryID
	count      Count
	occurredOn OccurredOn
}

// NewDoneLog constructs a DONELOG aggregate.
func NewDoneLog(
	id DoneLogID,
	title Title,
	trackID TrackID,
	categoryID CategoryID,
	count Count,
	occurredOn OccurredOn,
) (*DoneLog, error) {
	return &DoneLog{
		id:         id,
		title:      title,
		trackID:    trackID,
		categoryID: categoryID,
		count:      count,
		occurredOn: occurredOn,
	}, nil
}

// Update overwrites mutable fields of DONELOG.
func (d *DoneLog) Update(title Title, categoryID CategoryID, count Count, occurredOn OccurredOn) {
	d.title = title
	d.categoryID = categoryID
	d.count = count
	d.occurredOn = occurredOn
}

// Title returns the current title.
func (d *DoneLog) Title() Title {
	return d.title
}

// TrackID returns the Track identifier.
func (d *DoneLog) TrackID() TrackID {
	return d.trackID
}

// CategoryID returns the Category identifier.
func (d *DoneLog) CategoryID() CategoryID {
	return d.categoryID
}

// Count returns the number of things done.
func (d *DoneLog) Count() Count {
	return d.count
}

// OccurredOn exposes when the DONELOG happened.
func (d *DoneLog) OccurredOn() OccurredOn {
	return d.occurredOn
}

// ID returns the aggregate identifier.
func (d *DoneLog) ID() DoneLogID {
	return d.id
}
