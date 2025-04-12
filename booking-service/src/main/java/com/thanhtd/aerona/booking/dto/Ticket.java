package com.thanhtd.aerona.booking.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;
import lombok.Setter;

import java.util.Date;
import java.util.List;

@Getter
@Setter
public class Ticket {
    @JsonProperty("ticket_id")
    private String ticketId;

    @JsonProperty("pnr")
    private String PNR;

    @JsonProperty("booking_id")
    private String bookingId;

    @JsonProperty("flight_id")
    private String flightId;

    private List<PassengerInfo> passengers;

    @JsonProperty("ticket_number")
    private String ticketNumber;

    private String status;

    @JsonProperty("issued_at")
    private Date issuedAt;

    @JsonProperty("cancelled_at")
    private Date cancelledAt;

    @JsonProperty("pdf_url")
    private String pdfUrl;

    @JsonProperty("qr_code")
    private String qrCode;

    @JsonProperty("created_at")
    private Date createdAt;

    @JsonProperty("updated_at")
    private Date updatedAt;
}
