# Contributing to Castra

## How to Release

Castra uses GitHub Actions to automatically build and release binaries.

1.  **Commit your changes:** Ensure `main` is up to date.
2.  **Tag the release:**
    ```bash
    git tag v1.0.0
    ```
3.  **Push the tag:**
    ```bash
    git push origin v1.0.0
    ```

GitHub Actions will detect the tag, build `castra-mac`, `castra-linux`, and `castra-windows.exe`, and create a new Release on the GitHub repository page with these files attached.Users can then download the binary for their platform directly from the "Releases" section.
