package com.thanhtd.aerona.booking.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

import java.util.Date;

@Data
public class AirportLocation {
    @JsonProperty("airport_name")
    private String airportName;

    @JsonProperty("airport_code")
    private String airportCode;

    private String  city;

    private String country;

    private Date scheduled;
}
