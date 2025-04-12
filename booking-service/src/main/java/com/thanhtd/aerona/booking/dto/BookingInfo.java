package com.thanhtd.aerona.booking.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.List;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class BookingInfo {
    @JsonProperty("booking_id")
    private String bookingId;

    private String pnr;

    @JsonProperty("flight_id")
    private String flightId;

    private Flight flight;

    @JsonProperty("num_of_passengers")
    private Integer numOfPassengers;

    @JsonProperty("user_id")
    private String userId;

    private ContactInfo contact;

    private List<PassengerInfo> passengers;

    @JsonProperty("total_price")
    private Float totalPrice;

    private String currency;

    private String status;
}
