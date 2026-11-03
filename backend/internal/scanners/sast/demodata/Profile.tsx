type ProfileProps = { bio: string }

export function Profile({ bio }: ProfileProps) {
  const node = document.getElementById("bio")
  if (node) {
    node.innerHTML = bio
  }
  return <div dangerouslySetInnerHTML={{ __html: bio }} />
}
