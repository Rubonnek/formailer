# Getting Started with Standalone HBS form handler

This standalone example and custom form handler is specifically tailored to work with hugo-bootstrap-theme by Razon Yang.

Write `layouts/partials/hooks/contact/form-field-end.html`:

```html
<input type="hidden" name="_form_name" value="Contact">
```

And point the contact endpoint to your formailer instance.
