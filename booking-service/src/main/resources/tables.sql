create table booking (
    booking_id varchar(50) primary key,
    pnr varchar(25) unique not null,
    flight_id varchar(255) not null,
    user_id varchar(50),
    contact_id varchar(50),
    status varchar(15) not null,
    booked_at timestamp default current_timestamp,
    expires_at timestamp not null,
    confirmed_at timestamp,
    cancelled_at timestamp,
    total_price float not null,
    currency varchar(25) not null,
    num_of_passengers integer not null,
    updated_at timestamp
);

create table contact (
    contact_id varchar(50) primary key,
    booking_id varchar(50) not null,
    first_name varchar(75) not null,
    last_name varchar(75) not null,
    email varchar(75) not null,
    phone varchar(25) not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp,
    status varchar(15) not null,
    foreign key (booking_id) references booking(booking_id)
);

create table passenger (
    passenger_id varchar(50) primary key,
    booking_id varchar(50) not null,
    first_name varchar(75) not null,
    last_name varchar(75) not null,
    day_of_birth timestamp not null,
    nationality varchar(75) not null,
    passport_number varchar(74) not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp,
    status varchar(15) not null
);

create table booking_payment (
    id varchar(50) primary key,
    booking_id varchar(50) not null,
    payment_id varchar(50) not null,
    amount float not null,
    paid_at timestamp default current_timestamp,
    foreign key (booking_id) references booking(booking_id)
);