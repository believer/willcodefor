import type { Post, Prisma } from '@prisma/client'
import { prisma } from '~/db.server'

export function getPost(slug?: Post['slug']) {
  return prisma.post.findFirst({
    where: { OR: [{ slug }, { longSlug: slug }] },
    include: {
      _count: { select: { postViews: true } },
    },
  })
}

export type LatestTilPosts = Array<
  Pick<Post, 'id' | 'title' | 'tilId' | 'createdAt' | 'slug' | 'updatedAt'> & {
    _count: { postViews: number }
  }
>

export function getLatestTil({
  orderBy,
  take,
}: Pick<Prisma.PostFindManyArgs, 'orderBy' | 'take'>): Promise<LatestTilPosts> {
  return prisma.post.findMany({
    select: {
      _count: { select: { postViews: true } },
      tilId: true,
      title: true,
      id: true,
      slug: true,
      createdAt: true,
      updatedAt: true,
    },
    take,
    orderBy: orderBy ? orderBy : { createdAt: 'desc' },
  })
}

export function postSearch(query: string): Promise<LatestTilPosts> {
  return prisma.post.findMany({
    select: {
      _count: { select: { postViews: true } },
      tilId: true,
      title: true,
      id: true,
      slug: true,
      createdAt: true,
      updatedAt: true,
    },
    where: {
      OR: [
        {
          title: {
            mode: 'insensitive',
            search: decodeURI(query).replace(/\s/g, ' & '),
          },
        },
        {
          body: {
            mode: 'insensitive',
            search: decodeURI(query).replace(/\s/g, ' & '),
          },
        },
      ],
    },
    orderBy: { createdAt: 'desc' },
  })
}

export function getPosts() {
  return prisma.post.findMany({
    orderBy: { createdAt: 'desc' },
  })
}
