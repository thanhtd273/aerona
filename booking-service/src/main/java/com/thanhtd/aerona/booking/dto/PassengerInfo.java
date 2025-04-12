package com.thanhtd.aerona.booking.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.util.ObjectUtils;

import java.util.Date;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class PassengerInfo {
    @JsonProperty("first_name")
    private String firstName;

    @JsonProperty("last_name")
    private String lastName;

    @JsonProperty("day_of_birth")
    private Date dayOfBirth;

    private String nationality;

    @JsonProperty("passport_number")
    private String passportNumber;

    public boolean areFieldsFilled() {
        return !ObjectUtils.isEmpty(firstName) && !ObjectUtils.isEmpty(lastName) && !ObjectUtils.isEmpty(dayOfBirth)
                && !ObjectUtils.isEmpty(nationality) && !ObjectUtils.isEmpty(passportNumber);
    }
}
