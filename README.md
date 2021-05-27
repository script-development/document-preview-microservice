# Generate document previews micro service

A micro service for generating Generate document previews.
The most common document and image formats are supported!

## `/api/preview`

- Method: `POST`
- Body: `Multipart form`
- Auth: Bearer
- Supported file types
  - image/**png**, image/**jpeg**, image/**gif**, image/**webp**
  - application/**pdf**, _application/x-pdf, application/x-bzpdf, application/x-gzpdf_ (requires [pdftocairo from poppler](https://repology.org/project/poppler/versions) OR [soffice from openoffice](https://repology.org/project/openoffice/versions))
  - word, powerpoint and excel like files supported by openoffice (requires [soffice from openoffice](https://repology.org/project/openoffice/versions))

### Required fields:

- `document`: File
- `height`: Number
- `width`: Number

### Response:

- `200`: Image of content type `image/webp`
- `500`: Error string of type `text/plain`
- `401`: Unauthorized of type `text/plain`

## Env:

- `$KEY` **Required**: The Bearer key
- `$SOFFICE_PATH`: If the soffice binary is not in the path you can set it in here
- `$PORT`: Change the default port number (3030)
