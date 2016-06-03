# Manifest Format Specification

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be
interpreted as described in RFC 2119.

1. Files following the specification MUST contain a single JSON object.

2. It MUST NOT be the case that a key in that object is the prefix of any other
   key in that object.

3. That object MUST contain a key for each dependency.

4. The key MUST describe the path where the code from the dependency should be
   placed e.g. `x_limits_checker_libs`. The path MUST be a relative path.

5. The value for each key MUST be a JSON object.

6. That object MUST contain a key named "vcs".

7. The value of the key MUST describe the VCS (version control system) used to
   obtain the dependency.

8. If the dependency is obtained via SVN, the value MUST be "svn".

9. If the dependency is obtained via Git, then the value MUST be "git".

10. If the value of "vcs" is "svn" then:

    a. A key "url" MUST be present. It's value MUST contain the SVN URL from
    where the dependency can be checked out from.

    b. A key "rev" MAY be present. If present, it's value MUST contain the
    revision at which to obtain the dependency. If the key is not present,
    then the latest revision is assumed.

11. If the value of "vcs" is "git" then:

    a. A key "url" MUST be present. Its value MUST contain the URL from where
    the dependency's repository can be cloned from.

    b. A key "ref" MUST be present. Its value MUST contain a string that Git
    knows how to checkout. It could be a branch, a tag, a SHA1 hash, the output
    of git describe etc.

    c. A key "dir" MUST be present. Its value MUST indicate the directory
    inside the cloned repo that contains the dependency (if the value is an
    empty string, then this means the dependency is the contents of the whole
    repository). Its value MUST NOT be an absolute path or a Windows style path.

12. Other keys SHOULD NOT be present.

13. Files following the specification SHOULD reside in the root directory of
    the repository the dependencies are for.

