# Generate document previews micro service

Generate image previews

## `/api/preview`

- Method: `POST`
- Body: `Multipart form`
- Supported file types
  - image/**png**, image/**jpeg**, image/**gif**, image/**webp**
  - application/**pdf**, _application/x-pdf, application/x-bzpdf, application/x-gzpdf_ (requires [poppler](https://repology.org/project/poppler/versions))

### Required fields:

- `document`: File
- `height`: Number
- `width`: Number

### Response:

- `200`: image of content type `image/webp`
- `500`: error string of type `text/plain`

## Note on container:

The container uses the Microsoft TrueType core fonts that have a proprietary license
