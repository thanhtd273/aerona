package com.thanhtd.aerona.user.model;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import jakarta.persistence.Column;
import lombok.Data;

import java.io.Serializable;

@Entity
@Table(name = "user_role")
@Data
public class UserRole implements Serializable {

    @Id
    @Column(name = "id")
    private String id;

    @Column(name = "user_id")
    private String userId;

    @Column(name = "role_id")
    private String roleId;
}
