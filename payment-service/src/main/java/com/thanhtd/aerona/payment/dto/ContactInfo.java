package com.thanhtd.aerona.payment.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;
import org.springframework.util.ObjectUtils;

@AllArgsConstructor
@Getter
@Setter
public class ContactInfo {
    @JsonProperty("first_name")
    private String firstName;

    @JsonProperty("last_name")
    private String lastName;

    private String phone;

    private String email;

    public boolean areAllFieldsFilled() {
        return !ObjectUtils.isEmpty(firstName) && !ObjectUtils.isEmpty(lastName)
                && !ObjectUtils.isEmpty(email) && !ObjectUtils.isEmpty(phone);
    }
}
