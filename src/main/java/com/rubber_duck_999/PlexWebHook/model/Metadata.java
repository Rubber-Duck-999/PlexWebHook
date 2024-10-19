package com.rubber_duck_999.PlexWebHook.model;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class Metadata {
    private String librarySectionType;
    private String type;
    private String title;
    private String librarySectionTitle;
    private String grandparentTitle;
    private String parentTitle;
    private String contentRating;
    private String summary;
}