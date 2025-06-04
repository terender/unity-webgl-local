# unity-webgl-local
Convert all the datas in a Unity WebGL Build to js format. So it can be loaded and played from local files without runing a local web server.

# Dependency
- No required if using pre-compiled binary.
- Golang 1.20+ is needed if run or build from source code. 

# Usage
 - Make sure you are using customize <a href="https://docs.unity3d.com/2022.3/Documentation/Manual/webgl-templates.html">WebGL Template</a> in Unity WebGL Player Settings.
 - Copy the fetch.js into your WebGL Template directory. For example:
    ```
    <Your Unity Project>/Assets/WebGLTemplates/<Your Template>/TemplateData/fetch.js
    ```
   Edit the `index.html` file in your WebGL Template, add the following line to the `<head>` element:
   ```
   <script src ="TemplateData/fetch.js"></script>
   ```
   Make sure the path mathes.
 - Build player from Unity Editor, set the output directory name to `WebGL`. The `Compression Format` **MUST** be `Disabled`.
 - Copy the `convert-webgl.go` to the parent directory of the player build `WebGL` in last step. Run `go run .` in that path.
   Or else using the pre-compiled binary (for example `convert-webgl.go.exe` for windows) to that path and run.
 - The output directory `WebGLLocal` will be generated. Open the `index.html` with any browser you like.
 