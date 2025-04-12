package com.thanhtd.aerona.user.model;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import jakarta.persistence.Column;
import lombok.Data;

import java.util.Date;

@Table(name = "_action")
@Entity
@Data
public class Action {
    @Id
    @Column(name = "action_id")
    private String actionId;

    @Column(name = "name")
    private String name;

    @Column(name = "code")
    private Integer code;

    @Column(name = "description")
    private String description;

    @Column(name = "create_date")
    private Date createDate;

    @Column(name = "modified_date")
    private Date modifiedDate;

    @Column(name = "status")
    private String status;
}
