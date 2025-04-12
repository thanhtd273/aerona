package com.thanhtd.aerona.user.model;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.Data;

@Table(name = "role_permission")
@Entity
@Data
public class RolePermission {
    @Id
    @Column(name = "id")
    private String id;

    @Column(name = "role_id")
    private String roleId;

    @Column(name = "permission_id")
    private String permissionId;
}
