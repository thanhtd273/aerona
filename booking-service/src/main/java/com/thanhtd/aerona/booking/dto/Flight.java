package com.thanhtd.aerona.booking.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

import java.util.Date;
import java.util.List;

@Data
public class Flight {
    @JsonProperty("flight_id")
    private String flightId;

    @JsonProperty("flight_code")
    private String fightCode;

    @JsonProperty("departure")
    private AirportLocation departure;

    @JsonProperty("arrival")
    private AirportLocation arrival;

    @JsonProperty("duration")
    private Long duration;

    @JsonProperty("price")
    private Float price;

    private String currency;

    private String airline;

    @JsonProperty("seats_available")
    private Integer seatsAvailable;

    private Integer stops;

    @JsonProperty("stop_details")
    private List<StopDetail> stopDetails;

    @JsonProperty("airplane_type")
    private String airplaneType;

    @JsonProperty("flight_date")
    private String flightDate;

    @JsonProperty("created_at")
    private Date createdAt;

    @JsonProperty("updated_at")
    private Date updatedAt;

    private String status;
}
