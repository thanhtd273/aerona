package com.thanhtd.aerona.payment.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

import java.util.Date;

@Data
public class StopDetail {
    @JsonProperty("stop_id")
    private String stopId;

    @JsonProperty("airport_code")
    private String airportCode;

    @JsonProperty("departure_time")
    private Date departureTime;

    @JsonProperty("arrival_time")
    private Date arrivalTime;
}
