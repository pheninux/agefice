create table documents
(
    id      int auto_increment
        primary key,
    libelle varchar(50) not null
);

create table entreprises
(
    id   int auto_increment
        primary key,
    nom  varchar(50) null,
    code varchar(8)  null
);

create table formations
(
    id         int auto_increment
        primary key,
    created_at timestamp    null,
    intitule   varchar(255) null,
    date_debut timestamp    null,
    date_fin   timestamp    null,
    nbr_heures int          null,
    cout       double       null
);

create table personnes
(
    id              int auto_increment
        primary key,
    created_at      timestamp            null,
    nom             varchar(25)          not null,
    prenom          varchar(25)          null,
    age             int unsigned         null,
    date_naissance  date                 null,
    tel             varchar(11)          null,
    mail            varchar(255)         null,
    adresse         varchar(255)         null,
    entreprise_id   int                  null,
    flag_mail       tinyint(1) default 0 null,
    mfa             tinyint(1) default 0 null,
    nsocial         varchar(25)          null,
    status          varchar(20)          null,
    commentaire     varchar(250)         null,
    stop_mail       tinyint(1) default 0 null,
    prospection     tinyint(1) default 0 null,
    com_prospection varchar(250)         null,
    constraint personnes_nsocial_uindex
        unique (nsocial),
    constraint personnes_entreprise_id_entreprises_id_foreign
        foreign key (entreprise_id) references entreprises (id)
            on update cascade on delete cascade
);

create table personnes_documents
(
    personne_id int not null,
    document_id int not null,
    primary key (personne_id, document_id),
    constraint personnes_documents_document_id_documents_id_foreign
        foreign key (document_id) references documents (id)
            on update cascade on delete cascade,
    constraint personnes_documents_personne_id_personnes_id_foreign
        foreign key (personne_id) references personnes (id)
            on update cascade on delete cascade
);

create table personnes_formations
(
    personne_id  int not null,
    formation_id int not null,
    primary key (personne_id, formation_id),
    constraint personnes_formations_formation_id_formations_id_foreign
        foreign key (formation_id) references formations (id)
            on update cascade on delete cascade,
    constraint personnes_formations_personne_id_personnes_id_foreign
        foreign key (personne_id) references personnes (id)
            on update cascade on delete cascade
);

create table remboursements
(
    id           int auto_increment,
    formation_id int        not null,
    personne_id  int        not null,
    statut       tinyint(1) null,
    primary key (id, formation_id, personne_id),
    constraint remboursements_formation_id_formations_id_foreign
        foreign key (formation_id) references formations (id)
            on update cascade on delete cascade,
    constraint remboursements_personne_id_personnes_id_foreign
        foreign key (personne_id) references personnes (id)
            on update cascade on delete cascade
);

create table users
(
    id              int auto_increment
        primary key,
    name            varchar(255)         not null,
    email           varchar(255)         not null,
    hashed_password char(60)             not null,
    created_at      datetime             not null,
    active          tinyint(1) default 1 not null,
    constraint users_uc_email
        unique (email)
);

