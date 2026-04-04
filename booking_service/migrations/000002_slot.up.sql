CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL,
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('active', 'cancelled')),
    user_id UUID NOT NULL
);

CREATE INDEX idx_bookings_room_date ON bookings(room_id, date, status);
CREATE UNIQUE INDEX idx_unique_active_booking 
ON bookings(room_id, date, start_time) 
WHERE status = 'active';

CREATE INDEX idx_bookings_user_id ON bookings(user_id);