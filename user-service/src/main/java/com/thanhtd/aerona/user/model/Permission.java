package com.thanhtd.aerona.user.model;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import jakarta.persistence.Column;
import lombok.Data;

import java.io.Serializable;
import java.util.Date;

@Table(name = "_permission")
@Entity
@Data
public class Permission implements Serializable {
    @Id
    @Column(name = "permission_id")
    private String permissionId;

    @Column(name = "name")
    private String name;

    @Column(name = "description")
    private String description;

    @Column(name = "create_date")
    private Date createDate;

    @Column(name = "modified_date")
    private Date modifiedDate;

    @Column(name = "status")
    private String status;
}
