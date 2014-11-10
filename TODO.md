# TODO List

## Bugfixes
- [ ] Remove tags from document
- [ ] Add tags to document
- [ ] Quoted, multi-word tags

## Short-Term Features
- [ ] Allow documents to have multiple "pages"
- [ ] Support for cropping images (Jcrop)
- [ ] OCR support (Tesseract, maybe pyocr?)
  - Need to perform preprocessing on image first
    - [ ] Convert to grayscale
    - [ ] Rotation
    - [ ] De-skew
    - [ ] Border or corner detection
- [ ] Proper searching (using Whoosh)
  - Need to be able to re-index documents from source when schema changes
- [ ] Document-type specific metadata
  - [ ] Width/height of images
  - [ ] Number of pages in a PDF
  - [ ] Author information from various things
  - [ ] Metadata registration (i.e. handlers register for specific types, we
        create a proper Whoosh schema for each field, know how to display,
        etc. etc.)


## "Maybe in the Future" Features
- [ ] Folder/category support - for higher-level organizing
